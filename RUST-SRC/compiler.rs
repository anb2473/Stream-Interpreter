mod main;

use std::collections::HashMap;
use std::fs::{self, File};
use std::io::{Read, Write};
use std::path::{Path, PathBuf};
use std::process::exit;
use std::sync::{Arc, Mutex};
use std::{io, thread};
use std::time::{Duration, Instant};
use std::env;
use regex::Regex;
use serde_json::Value;

// Define print variables
const RESET: &str = "\x1b[0m";
const BOLD: &str = "\x1b[1m";
const ITALIC: &str = "\x1b[3m";

const FG_BRIGHT_RED: &str = "\x1b[91m";
const FG_BRIGHT_GREEN: &str = "\x1b[92m";
const FG_BRIGHT_YELLOW: &str = "\x1b[93m";
const FG_BRIGHT_BLUE: &str = "\x1b[94m";
const FG_BRIGHT_CYAN: &str = "\x1b[96m";

static START_TIME: Mutex<Option<Instant>> = Mutex::new(None);

fn start_timer() {
    let mut start_time = START_TIME.lock().unwrap();
    *start_time = Some(Instant::now());
}

fn get_elapsed_time() -> Option<Duration> {
    let start_time = START_TIME.lock().unwrap();
    start_time.map(|start| start.elapsed())
}

fn read_file(file_path: &str) -> Result<String, io::Error> {
    let path = Path::new(file_path);
    let mut file = File::open(path)?; // Open the file, propagate errors

    let mut contents = String::new();
    file.read_to_string(&mut contents)?; // Read contents, propagate errors

    Ok(contents) // Return the contents if successful
}

fn write_file(file_path: &str, contents: &str) -> Result<(), io::Error> {
    let path = Path::new(file_path);
    let mut file = File::create(path)?; // Create or truncate the file, propagate errors

    file.write_all(contents.as_bytes())?; // Write the contents, propagate errors

    Ok(()) // Return Ok if successful
}

fn get_basename(path: &str) -> Option<&str> {
    Path::new(path).file_name().and_then(|name| name.to_str())
}

fn has_lone_equals(text: &str) -> bool {
    let chars: Vec<char> = text.chars().collect();

    for i in 0..chars.len() {
        if chars[i] == '=' {
            let is_lone_equal = (i == 0 || chars[i - 1] != '=') &&
                (i + 1 == chars.len() || chars[i + 1] != '=');

            if is_lone_equal {
                return true; // Found a lone '=', return true immediately
            }
        }
    }
    false // No lone equals found
}

fn mkdir(path: &str) -> i32 {
    let path = Path::new(path);
    match fs::create_dir_all(path) {
        Ok(_) => 0, // Success
        Err(err) => {
            eprintln!("Error creating directory: {}", err);
            err.raw_os_error().unwrap_or(-1) // Return the OS error code or -1
        }
    }
}

struct App {
    path: String,
    output: Option<String>,
    syntactic: bool,
    ext: String,
    file: String,
    active_directory: String,
    inheritance_map: HashMap<String, String>,
    local: HashMap<String, String>,
}

impl App {
    fn new(path: &str, output: Option<&str>, syntactic: bool) -> Result<Self, String> {
        let path_obj = Path::new(path);
        let ext = path_obj.extension()
            .map_or("".to_string(), |ext| format!(".{}", ext.to_string_lossy()));

        // CHECK FOR VALID APPLICATION EXTENSION
        if ext != ".stream" && ext != ".sl" {
            println!("{}{:?}{}: {}{}Invalid file extension:{} {}'{}'{}",
                     FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                     BOLD, FG_BRIGHT_RED, RESET,
                     ITALIC, ext, RESET);
            exit(1);
        }

        // CHECK FILE PATH EXISTS
        if !path_obj.exists() {
            println!("{}{:?}{}: {}{}File not found:{} {}'{}'{}",
                     FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                     BOLD, FG_BRIGHT_RED, RESET,
                     ITALIC, path, RESET);
            exit(1);
        }

        // LOAD FILE CONTENTS
        let file = match read_file(path) {
            Ok(content) => content,
            Err(e) => return Err(format!("Failed to read file: {}", e)), // Convert io::Error to Error
        };

        // GET WORKING DIRECTORY OF THE PATH
        let active_directory = path_obj.parent()
            .unwrap_or_else(|| Path::new(""))
            .to_string_lossy()
            .to_string();

        // BUILD 'build' DIRECTORY
        let build_dir = format!("{}{}build", active_directory, std::path::MAIN_SEPARATOR);
        if mkdir(&build_dir) != 0 {
            exit(1);
        }

        // IF NO OUTPUT DIRECTORY IS SPECIFIED, GENERATE AN OUTPUT DIRECTORY
        let basename = path_obj.file_stem().unwrap().to_string_lossy();
        let output = output.map(|o| o.to_string()).unwrap_or_else(|| {
            format!("{}{}build{}build_{}{}{}",
                    active_directory, std::path::MAIN_SEPARATOR,
                    std::path::MAIN_SEPARATOR, basename,
                    std::path::MAIN_SEPARATOR, basename)
        });

        // BUILD SUB BUILD DIRECTORY
        let sub_build_dir = format!("{}{}build{}build_{}",
                                    active_directory, std::path::MAIN_SEPARATOR,
                                    std::path::MAIN_SEPARATOR, basename);
        if mkdir(&sub_build_dir) != 0 {
            exit(1);
        }

        // BUILD FUNCTION DIRECTORY
        let func_dir = format!("{}{}functions", sub_build_dir, std::path::MAIN_SEPARATOR);
        if mkdir(&func_dir) != 0 {
            exit(1);
        }

        Ok(App {
            path: path.to_string(),
            output: Some(output),
            syntactic,
            ext,
            file,
            active_directory,
            inheritance_map: HashMap::new(),
            local: HashMap::new(),
        })
    }

    fn run(&mut self) {
        if self.ext == ".stream" {  // CORRECT FILE FOR COMPILATION IF FILE IS A HIGH LEVEL STREAM APPLICATION
            self.file = self.correct_safety();

            println!("{}{:?}{}: {}Successfully corrected file safety{}",
                     FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                     FG_BRIGHT_GREEN, RESET);

            // WRITE STREAM LITE FILE
            write_file(&format!("{}.sl", self.output.as_ref().unwrap()), &self.file);
        }

        // GENERATE FUNCTIONAL STREAM LITE FILE
        self.file = self.generate_functional_file();

        println!("{}{:?}{}: {}Successfully generated functional stream lite{}",
                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                 FG_BRIGHT_GREEN, RESET);

        // WRITE FUNCTIONAL STREAM LITE FILE
        write_file(&format!("{}.fsl", self.output.as_ref().unwrap()), &self.file);

        println!("{}{:?}{}: {}Successfully finished compilation{}",
                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                 FG_BRIGHT_GREEN, RESET);
    }

    fn get_type(&self, var: &str, line_num: usize) -> String {
        let var = var.trim();
        if var.is_empty() {
            return "obj".to_string();
        }

        if var.contains('.') {
            if let Some(obj) = self.local.get(var.split('.').next().unwrap()) {
                return obj.rsplit('>').next().unwrap().trim().to_string();
            }
        }

        if let Some(t) = self.local.get(var) {
            return t.clone();
        }

        let operators = ['+', '-', '*', '/', '~', '$', '@', '&', '%'];
        if operators.iter().any(|&op| var.contains(op)) {
            let mut true_var = String::new();
            let mut i = 0;

            for c in var.chars() {
                if operators.contains(&c) {
                    break;
                } else {
                    true_var.push(c);
                }
                i += 1;
            }

            return self.get_type(&true_var, line_num);
        }

        match var.chars().next() {
            Some('\'') => {
                // LINE STARTS WITH CHAR TYPE QUOTE
                if var.replace('\'', "").len() > 1 {
                    println!("{}{:?}{}: {}{}\
                             Value Warning:{} {}'{}'{}  \
                             is not a single character, {}line {}{}",
                             FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                             BOLD, FG_BRIGHT_YELLOW, RESET,
                             ITALIC, var, RESET,
                             FG_BRIGHT_BLUE, line_num, RESET);
                }
                "char".to_string()
            },
            Some('\"') => "str".to_string(), // LINE STARTS WITH STRING TYPE QUOTE
            Some('[') => {
                // LINE STARTS WITH LIST DECLARATION
                let first_item = var[1..].split(',').next().unwrap_or("");
                format!("list > {}", self.get_type(first_item, line_num))
            },
            Some('{') => {
                // LINE STARTS WITH DICTIONARY DECLARATION
                let first_item = var[1..].split(',').next().unwrap_or("");
                let key_value: Vec<&str> = first_item.split(':').collect();
                if key_value.len() >= 2 {
                    format!("dict > {} > {}",
                            self.get_type(key_value[0], line_num),
                            self.get_type(key_value[1], line_num))
                } else {
                    "dict > obj > obj".to_string()
                }
            },
            _ => {
                if var.starts_with("if") || var == "true" || var == "false" {
                    // LINE IS BOOLEAN OR BOOLEAN EXPRESSION
                    "bool".to_string()
                } else if var.chars().all(|c| c.is_digit(10)) {
                    // LINE IS INTEGER
                    "int".to_string()
                } else if var.parse::<f64>().is_ok() {
                    "float".to_string()
                } else {
                    "obj".to_string()
                }
            }
        }
    }

    fn correct_safety(&mut self) -> String {
        let mut new_file = String::new();
        let mut namespace: Option<String> = None;

        let mut i = 0;
        for line in self.file.split('\n') {
            i += 1;
            // SPLIT MULTI ACTION LINES INTO SUB LINES
            let split_line: Vec<&str> = line.split(';').collect();

            for sub_line in split_line {
                // REMOVE COMMENTS AND REDUNDANT WHITESPACES
                let mut new_line = sub_line.split('#').next().unwrap().trim().to_string();

                if new_line.starts_with("if") || new_line.starts_with("elif") ||
                    new_line.starts_with("while") || new_line.starts_with("else") {
                    if (new_line.starts_with("w") && new_line.contains('=')) ||
                        !new_line.starts_with("w") {
                        let upper_checks = ["and", "or"];
                        let tags = ["not"];

                        let mut ret = String::new();
                        let split_line: Vec<&str>;

                        if new_line.starts_with("i") {
                            split_line = new_line[3..].split('{').next().unwrap().split(' ').collect();
                            ret = "if ".to_string();
                        } else if new_line.starts_with("w") {
                            let parts: Vec<&str> = new_line[6..].split('{').next().unwrap().split('=').collect();
                            if parts.len() >= 2 {
                                let condition_parts: Vec<&str> = parts[1][3..].split(' ').collect();
                                split_line = condition_parts;
                                ret = format!("while {} = if ", parts[0].trim());
                            } else {
                                split_line = Vec::new();
                                ret = "while ".to_string();
                            }
                        } else if new_line.starts_with("els") {
                            split_line = new_line[5..].split('{').next().unwrap().split(' ').collect();
                            ret = "else ".to_string();
                        } else if new_line.starts_with("eli") {
                            split_line = new_line[5..].split('{').next().unwrap().split(' ').collect();
                            ret = "elif ".to_string();
                        } else {
                            split_line = Vec::new();
                            ret = "".to_string();
                        }

                        for part in split_line {
                            if !part.is_empty() {
                                if !part.contains('=') && !part.contains('<') &&
                                    !part.contains('>') && !upper_checks.contains(&part) &&
                                    !tags.contains(&part) && !part.contains(':') {
                                    ret += &format!("{}:{} ", part.trim(), self.get_type(part, i));
                                    println!("{}{:?}{}: {}{}\
                                             Typeless conditional object: '{}', {} {}'{}'{}\
                                             , {}line {}{}",
                                             FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                             BOLD, FG_BRIGHT_YELLOW, part, RESET,
                                             ITALIC, line, RESET,
                                             FG_BRIGHT_BLUE, i, RESET);
                                } else {
                                    ret += &format!("{} ", part.trim());
                                }
                            }
                        }

                        new_line = format!("{}{}", ret, "{");
                    }

                    if !new_line.contains('{') {
                        new_line += "{";
                        println!("{}{:?}{}: {}{}No end braces:{} {}'{}'{}, {}line {}{}",
                                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                 BOLD, FG_BRIGHT_YELLOW, RESET,
                                 ITALIC, line, RESET,
                                 FG_BRIGHT_BLUE, i, RESET);
                    }
                } else if new_line.starts_with("def") {
                    // ENSURE LINE CORRECTLY ENDS WITH AN ENDING TAG
                    if !new_line.contains('{') {
                        new_line += "{";
                        println!("{}{:?}{}: {}{}No end braces:{} {}'{}'{}, {}line {}{}",
                                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                 BOLD, FG_BRIGHT_YELLOW, RESET,
                                 ITALIC, line, RESET,
                                 FG_BRIGHT_BLUE, i, RESET);
                    }

                    if !new_line.contains('(') {
                        new_line = format!("{} ()", new_line.split('{').next().unwrap().trim());
                    }

                    let params_section = new_line.split('(').nth(1).unwrap_or("").split(')').next().unwrap_or("");
                    for param_pair in params_section.split(',') {
                        if !param_pair.is_empty() && param_pair.contains(':') {
                            let param_parts: Vec<&str> = param_pair.split(':').collect();
                            if param_parts.len() >= 2 {
                                self.local.insert(param_parts[0].trim().to_string(), param_parts[1].trim().to_string());
                            }
                        }
                    }

                    // ENSURE LINE CONTAINS A RETURN TYPE
                    if !new_line.contains("->") {
                        // INSERT VOID RETURN
                        new_line = format!("{} -> void ({}",
                                           new_line.split('(').next().unwrap().trim(),
                                           new_line.split('(').nth(1).unwrap_or(""));
                        println!("{}{:?}{}: {}{}\
                                 No function return notation:{} {}'{}'{}\
                                 , {}line {}{}",
                                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                 BOLD, FG_BRIGHT_YELLOW, RESET,
                                 ITALIC, line, RESET,
                                 FG_BRIGHT_BLUE, i, RESET);
                    }

                    // UPDATE NAME WITH NAMESPACE TAG
                    if let Some(ns) = &namespace {
                        new_line = format!("def {}.{}", ns, &line[4..]);
                    }
                } else if new_line.starts_with("return") {
                    if new_line.split(' ').count() != 2 {
                        println!("{}{:?}{}: {}{}\
                                 Return without statement:{} {}'{}'{}\
                                 , {}line {}{}",
                                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                 BOLD, FG_BRIGHT_YELLOW, RESET,
                                 ITALIC, line, RESET,
                                 FG_BRIGHT_BLUE, i, RESET);
                        new_line = format!("{} void", new_line.trim());
                    }

                    if new_line.contains(')') && new_line.contains('(') && !new_line.contains("exec") {
                        new_line = format!("return exec {}", &new_line[7..]);
                        println!("{}{:?}{}: {}{}\
                                 No exec tag:{} {}'{}'{}\
                                 , {}line {}{}",
                                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                 BOLD, FG_BRIGHT_YELLOW, RESET,
                                 ITALIC, line, RESET,
                                 FG_BRIGHT_BLUE, i, RESET);
                    }
                } else if new_line.starts_with("namespace") {
                    namespace = Some(new_line[10..].to_string());
                    if &namespace.as_ref().unwrap() == &"end" {
                        namespace = None;
                    }
                    new_line = String::new();
                } else if new_line.starts_with("end") {
                    namespace = None;
                    new_line = String::new();
                }

                // IF LINE CONTAINS = THAT IS NOT A PART OF A ==
                let re = Regex::new(r"([^=]|^)=([^=]|$)").unwrap();
                if re.is_match(&new_line) && !new_line.starts_with("while") {
                    // IF LINE DOES NOT HAVE AN EXEC TAG
                    if new_line.trim().ends_with(')') && new_line.contains('(') && !new_line.contains("exec") {
                        let parts: Vec<&str> = new_line.split('=').collect();
                        new_line = format!("{} = exec {}", parts[0].trim(), parts[1..].join("=").trim());
                        println!("{}{:?}{}: {}{}\
                                 No exec tag:{} {}'{}'{}\
                                 , {}line {}{}",
                                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                 BOLD, FG_BRIGHT_YELLOW, RESET,
                                 ITALIC, line, RESET,
                                 FG_BRIGHT_BLUE, i, RESET);
                    }

                    // IF LINE DOES NOT HAVE A LET DECLARATION INSERT LET
                    if !new_line.starts_with("let") {
                        new_line = format!("let {}", new_line);
                        println!("{}{:?}{}: {}{}\
                                 No let tag:{} {}'{}'{}, {}line {}{}",
                                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                 BOLD, FG_BRIGHT_YELLOW, RESET,
                                 ITALIC, line, RESET,
                                 FG_BRIGHT_BLUE, i, RESET);
                    }

                    // IF LINE DOES NOT HAVE A TYPE ANNOTATION GENERATE TYPE NOTATION
                    if !new_line.split('=').next().unwrap().contains(':') {
                        let split_eq: Vec<&str> = new_line.splitn(2, '=').collect();
                        let var = split_eq[1].trim();
                        let hand = split_eq[0];

                        let type_declaration = self.get_type(var, i);

                        new_line = format!("{}: {} = {}", hand.trim_end(), type_declaration, var);

                        println!("{}{:?}{}: {}{}\
                                 No variable type notation:{} {}'{}'{}\
                                 , {}line {}{}",
                                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                 BOLD, FG_BRIGHT_YELLOW, RESET,
                                 ITALIC, line, RESET,
                                 FG_BRIGHT_BLUE, i, RESET);
                    }
                }

                if new_line.contains(')') && new_line.contains('(') &&
                    !new_line.contains("exec") && !new_line.starts_with("let") &&
                    !new_line.starts_with("prototype") {
                    new_line = format!("exec {}", new_line);
                    println!("{}{:?}{}: {}{}\
                             No exec tag:{} {}'{}'{}\
                             , {}line {}{}",
                             FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                             BOLD, FG_BRIGHT_YELLOW, RESET,
                             ITALIC, line, RESET,
                             FG_BRIGHT_BLUE, i, RESET);
                }

                // IF LINE EXECUTES FUNCTION BUT DOES NOT CONTAIN AN EQUALS STATEMENT INSERT VOID STATEMENT
                if new_line.contains("exec") && !has_lone_equals(&new_line) && !line.starts_with("return") {
                    new_line = format!("let _: obj = {}", new_line);
                    println!("{}{:?}{}: {}{}\
                             No let statement:{} {}'{}'{}\
                             , {}line {}{}",
                             FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                             BOLD, FG_BRIGHT_YELLOW, RESET,
                             ITALIC, line, RESET,
                             FG_BRIGHT_BLUE, i, RESET);
                }

                if new_line.contains("exec") {
                    let parts: Vec<&str> = new_line.splitn(2, '(').collect();
                    let params_part = if parts.len() > 1 {
                        let end_parts: Vec<&str> = parts[1].rsplitn(2, ')').collect();
                        if end_parts.len() > 1 {
                            Some((end_parts[1], end_parts[0]))
                        } else {
                            None
                        }
                    } else {
                        None
                    };

                    if let Some((params_str, end_part)) = params_part {
                        let func_params: Vec<&str> = params_str
                            .split(',')
                            .map(|s| s.trim())
                            .filter(|s| !s.is_empty())
                            .collect();

                        let mut ret = format!("{}(", parts[0]);
                        for (index, param) in func_params.iter().enumerate() {
                            let mut modified_param = param.to_string();
                            if param.contains('(') && param.contains(')') {
                                modified_param = format!("exec {}", param.trim());
                            }

                            if index != func_params.len() - 1 {
                                ret += &format!("{},", modified_param);
                            } else {
                                ret += &modified_param;
                            }
                        }

                        new_line = format!("{}){}", ret, end_part);
                        println!("{}{:?}{}: {}{}\
                                 No exec tag:{} {}'{}'{}\
                                 , {}line {}{}",
                                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                 BOLD, FG_BRIGHT_YELLOW, RESET,
                                 ITALIC, line, RESET,
                                 FG_BRIGHT_BLUE, i, RESET);
                    }
                }

                if new_line.starts_with("let") {
                    // IF LINE HAS A DECLARATION BUT NOT A = STATEMENT INSERT VOID
                    if !has_lone_equals(&new_line) {
                        new_line = format!("{} = void", new_line.trim_end());
                        println!("{}{:?}{}: {}{}\
                                 No equals statement:{} {}'{}'{}\
                                 , {}line {}{}",
                                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                 BOLD, FG_BRIGHT_YELLOW, RESET,
                                 ITALIC, line, RESET,
                                 FG_BRIGHT_BLUE, i, RESET);
                    } else if new_line.trim().ends_with('=') {
                        new_line = format!("{} void", new_line.trim());
                        println!("{}{:?}{}: {}{}\
                                 No equals statement:{} {}'{}'{}\
                                 , {}line {}{}",
                                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                 BOLD, FG_BRIGHT_YELLOW, RESET,
                                 ITALIC, line, RESET,
                                 FG_BRIGHT_BLUE, i, RESET);
                    }

                    if line.contains('.') {
                        let parts: Vec<&str> = new_line.split('=').collect();
                        if parts.len() > 1 {
                            let right_parts: Vec<&str> = parts[1].split('.').collect();
                            if right_parts.len() > 1 {
                                let split_line = right_parts[1].trim();
                                if !split_line.contains(':') {
                                    new_line = format!("{}:{}", new_line.trim_end(), self.get_type(split_line, i));
                                }
                            }
                        }
                    }

                    if new_line.contains("if") {
                        let upper_checks = ["and", "or"];
                        let tags = ["not"];

                        let parts: Vec<&str> = new_line.splitn(2, '=').collect();
                        if parts.len() > 1 && parts[1].len() > 3 && parts[1].starts_with(" if ") {
                            let split_line: Vec<&str> = parts[1][3..].trim().split(' ').collect();
                            let mut ret = format!("let {} = if ", parts[0][4..].trim());

                            for part in split_line {
                                if !part.is_empty() {
                                    if !part.contains('=') && !upper_checks.contains(&part) &&
                                        !tags.contains(&part) && !part.contains(':') {
                                        ret += &format!("{}:{} ", part.trim(), self.get_type(part, i));
                                        println!("{}{:?}{}: {}{}\
                                                 Typeless conditional object: '{}', {} {}'{}'{}\
                                                 , {}line {}{}",
                                                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                                 BOLD, FG_BRIGHT_YELLOW, part, RESET,
                                                 ITALIC, line, RESET,
                                                 FG_BRIGHT_BLUE, i, RESET);
                                    } else {
                                        ret += &format!("{} ", part.trim());
                                    }
                                }
                            }

                            new_line = ret;
                        }
                    }

                    if let Some(last_char) = line.chars().last() {
                        if last_char == '*' || last_char == '/' || last_char == '+' || last_char == '-' {
                            println!("{}{:?}{}: {}{}\
                                     Unclosed operation:{} {}'{}'{}\
                                     , {}line {}{}",
                                     FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                     BOLD, FG_BRIGHT_YELLOW, RESET,
                                     ITALIC, line, RESET,
                                     FG_BRIGHT_BLUE, i, RESET);
                        }
                    }

                    // INSERT NAMESPACE TAG IN VARIABLE NAME
                    if let Some(ns) = &namespace {
                        new_line = format!("let {}.{}", ns, &new_line[4..]);
                        println!("{}{:?}{}: {}{}\
                                 No let tag:{} {}'{}'{}\
                                 , {}line {}{}",
                                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                                 BOLD, FG_BRIGHT_YELLOW, RESET,
                                 ITALIC, line, RESET,
                                 FG_BRIGHT_BLUE, i, RESET);
                    }

                    // Store variable types in local map
                    let var_parts: Vec<&str> = new_line[4..].split(':').collect();
                    if var_parts.len() > 1 {
                        let var_name = var_parts[0].trim();
                        let type_parts: Vec<&str> = var_parts[1].split('=').collect();
                        if !type_parts.is_empty() {
                            self.local.insert(var_name.to_string(), type_parts[0].trim().to_string());
                        }
                    }
                }

                // INSERT LINE
                if !new_line.is_empty() {
                    new_file += &format!("{}\n", new_line);
                }
            }
        }

        new_file
    }

    fn generate_functional_file(&self) -> String {
        let mut new_file = String::new();
        let mut in_foo = false;
        let mut split_file: Vec<&str> = self.file.split('\n').collect();
        let mut index = 0;
        let mut in_while = false;
        let mut while_var = vec!["".to_string()];
        let mut while_jump_num = Vec::new();
        let mut general_depth = 0;
        let mut while_general_depth = Vec::new();

        // LOOP THROUGH LINES IN FILE
        while index < split_file.len() {
            // RETRIEVE LINE INFO
            let line = split_file[index];
            let mut new_line = line.to_string();

            // IF LINE IS AN IMPORT
            if line.starts_with("from") {
                let split_line: Vec<&str> = line[5..].split("import").collect();
                let directory = format!("{}/{}", self.active_directory, split_line[0].trim().trim_matches(|c| c == '"' || c == '\''));
                new_line = String::new();

                // CHECK IF DIRECTORY EXISTS
                if Path::new(&directory).exists() {
                    let functions: Vec<&str> = if split_line.len() > 1 {
                        split_line[1].split(',').map(|s| s.trim()).collect()
                    } else {
                        Vec::new()
                    };

                    // READ PACKAGE DATA
                    let pk_file_path = format!("{}/function_manager.spk", directory);
                    let mut pk_file = String::new();
                    if let Ok(mut file) = File::open(&pk_file_path) {
                        let _ = file.read_to_string(&mut pk_file);
                    }

                    let pk_lines: Vec<&str> = pk_file.split('\n').collect();

                    // SPLIT PACKAGE DATA BY FUNCTION
                    let mut split_pk_file = HashMap::new();
                    for pk_line in pk_lines {
                        if pk_line.contains("<%$>") {
                            let parts: Vec<&str> = pk_line.split("<%$>").collect();
                            split_pk_file.insert(parts[0].trim(), parts[1].trim());
                        }
                    }

                    for function in functions {
                        let stripped_function = function.trim();
                        if stripped_function == "*" {
                            for value in &split_pk_file {
                                new_line += &format!("prototype {} \"{}/{}\" {}\n",
                                    value.0.split(".").next().unwrap_or(""),
                                    directory,
                                    value.0,
                                    value.1
                                )
                            }
                        }
                        else {
                            new_line += &format!("prototype {} \"{}/{}\" {}\n",
                                function.split(".").next().unwrap_or(""),
                                directory,
                                function,
                                match split_pk_file.get(function) {
                                 Some(s) => format!("{}", s.to_string()), // Dereference twice and format.
                                 None => "None".to_string(), // Handle None case.
                                }
                            )
                        }
                    }
                }
                else {
                    println!(
                        "{}{:?}{}: {}{}No such directory:{} Unable to import functions at '{}/{}'",
                        FG_BRIGHT_CYAN,
                        get_elapsed_time(),
                        RESET,
                        BOLD,
                        FG_BRIGHT_RED,
                        RESET,
                        self.active_directory,
                        directory
                    );
                    if self.syntactic {
                        exit(1);
                    }
                }

                new_line = (new_line.trim_end()).to_string()
            }
            else if line.starts_with("def") {
                in_foo = true;
                let parts: Vec<&str> = line.split(' ').collect();
                let function_name = parts[4];
                let base_name = Path::new(&self.path)
                    .file_name()
                    .unwrap()
                    .to_str()
                    .unwrap()
                    .split('.')
                    .collect::<Vec<&str>>()[0];
                let dest_type = line.split("->").collect::<Vec<&str>>()[1].trim().split(' ').collect::<Vec<&str>>()[0];
                let new_line = format!("prototype {} \"{}\\build\\build_{}\\functions\\{}.fsl\" -> {} {}", function_name, self.active_directory, base_name, function_name, dest_type, if let Some(open_paren_index) = line.find('(') { let close_paren_index = line.find(')').unwrap(); let parameter_str = &line[open_paren_index + 1..close_paren_index]; if parameter_str.is_empty() { "()".to_string() } else { parameter_str.split(',').collect::<Vec<&str>>().join(", ").replace('\'', "").replace('[', "(").replace(']', ")") } } else { "()".to_string() });
                let mut depth = 0;
                let mut data: String = String::new();
                let mut jump_num = 1;
                for sub_line in split_file.iter().skip(index + 1) {
                    if sub_line.starts_with("if") || sub_line.starts_with("else") || sub_line.starts_with("elif") || sub_line.starts_with("while") {
                        depth += 1;
                    }
                    else if *sub_line == "}" {
                        depth -= 1;
                        if depth <= 0 {
                            break;
                        }
                    }
                    jump_num += 1;
                    data = format!("{}{}\n", data, sub_line);
                }
                index += jump_num;
                let base_name = Path::new(&self.path).file_name().and_then(|f| f.to_str()).map(|s| s.split('.').next().unwrap_or(s)).unwrap_or("");
                let function_name = line[4..].split(' ').next().unwrap_or("");
                let file_path = format!(
                    "{}\\build\\build_{}\\functions\\{}.sl",
                    self.active_directory, base_name, function_name
                );
                write_file(&file_path, &data);

                let mut sub_app: Result<App, String> = App::new(
                    &format!(
                        "{}\\build\\build_{}\\functions\\{}.sl",
                        self.active_directory,
                        {
                            let parts: Vec<&str> = self.path.split(std::path::MAIN_SEPARATOR).collect();
                            parts.last().map(|s| s.to_string()).unwrap_or_else(String::new)
                        },
                        line[4..].split(' ').next().unwrap_or("")
                    ),
                    Some(format!(
                        "{}\\build\\build_{}\\functions\\{}",
                        self.active_directory,
                        {
                            let parts: Vec<&str> = self.path.split(std::path::MAIN_SEPARATOR).collect();
                            parts.last().map(|s| s.to_string()).unwrap_or_else(String::new)
                        },
                        line[4..].split(' ').next().unwrap_or("")
                    ).as_str()),
                    self.syntactic
                );
                match sub_app {
                    Ok(mut app) => {
                        app.run(); // Correct: call run() on the App instance
                    }
                    Err(err) => {
                        eprintln!("Error creating App: {}", err);
                    }
                }
            }
            else if line.starts_with("if") {
                if !while_general_depth.is_empty() {
                    *while_general_depth.last_mut().unwrap() += 1;
                }
                general_depth += 1;
                new_file += &format!("let check: bool = if {}", line[3..line.len() - 1].trim());
                let mut depth = 0;
                let mut jump_num = 0;
                for sub_line in &mut split_file[index..] {
                    if sub_line.starts_with("if") || sub_line.starts_with("else") || sub_line.starts_with("elif") || sub_line.starts_with("while") {
                        depth += 1;
                    }
                    else if *sub_line == "}" {
                        depth -= 1;
                        if depth <= 0 {
                            break;
                        }
                    }
                    jump_num += 1;
                }
                new_line = format!("do check {}", jump_num - 1);
            }
            else if line.starts_with("else") {
                if !while_general_depth.is_empty() {
                    *while_general_depth.last_mut().unwrap() += 1;
                }
                general_depth += 1;
                new_file += "let check: bool = if not check:bool == true:bool";
                let mut depth = 0;
                let mut jump_num = 0;
                for sub_line in &mut split_file[index..] {
                    if sub_line.starts_with("if") || sub_line.starts_with("else") || sub_line.starts_with("elif") || sub_line.starts_with("while") {
                        depth += 1;
                    }
                    else if *sub_line == "}" {
                        depth -= 1;
                        if depth <= 0 {
                            break;
                        }
                    }
                    jump_num += 1;
                }
                new_line = format!("do check {}", jump_num - 1);
            }
            else if line.starts_with("elif") {
                if !while_general_depth.is_empty() {
                    *while_general_depth.last_mut().unwrap() += 1;
                }
                general_depth += 1;
                new_file += &format!("let check: bool = if not check:bool == true:bool and {}", line[5..line.len() - 1].trim());
                let mut depth = 0;
                let mut jump_num = 0;
                for sub_line in &mut split_file[index..] {
                    if sub_line.starts_with("if") || sub_line.starts_with("else") || sub_line.starts_with("elif") || sub_line.starts_with("while") {
                        depth += 1;
                    }
                    else if *sub_line == "}" {
                        depth -= 1;
                        if depth <= 0 {
                            break;
                        }
                    }
                    jump_num += 1;
                }
                new_line = format!("do check {}", jump_num - 1);
            }
            else if line.starts_with("while") {
                let mut depth = 0;
                let mut jump_num = 0;
                for sub_line in &mut split_file[index..] {
                    if sub_line.starts_with("if") || sub_line.starts_with("else") || sub_line.starts_with("elif") || sub_line.starts_with("while") {
                        depth += 1;
                    }
                    if sub_line.starts_with("while") && *sub_line != line {
                        jump_num += 2;
                    }
                    else if *sub_line == "}" {
                        depth -= 1;
                        if depth <= 0 {
                            break;
                        }
                    }
                    jump_num += 1;
                }
                let val_string: String = new_line[6..].split("{").next().unwrap_or("").trim().to_string();
                let mut val: &str = val_string.as_str();
                if val.contains("=") {
                    if val.starts_with("let") {
                        val = &val[4..];
                    }
                    new_file += &format!(
                        "let {}: bool = {}\n",
                        val.split('=').next().unwrap().split(':').next().unwrap().trim(),
                        val.split('=').nth(1).unwrap().trim()
                    );
                    let val = val.split('=').next().unwrap().split(':').next().unwrap().trim();
                }
                new_line = format!("do {} {}", val, jump_num - 1);

                while_jump_num.push(jump_num - 1);
                while_general_depth.push(1);
                in_while = true;
                while_var.push(val.to_string());
            }
            else if line == "}" {
                general_depth -= 1;
                if !while_general_depth.is_empty() {
                    *while_general_depth.last_mut().unwrap() -= 1;
                }
                if in_foo {
                    if index > 0 && !split_file[index - 1].starts_with("return") {
                        new_file.push_str("return void\n");
                    }
                }
                else if in_while && *while_general_depth.last().unwrap() == 0 {
                    while_general_depth.pop();
                    let check_var = while_var.pop().unwrap(); // Get the variable name
                    let jump_num = -while_jump_num.pop().unwrap(); // Get the jump number
                    new_line = format!("let check: bool = if not {}:bool == true:bool\ndo check {}", check_var, jump_num);
                }
                else {
                    new_line = String::new();
                }
                in_foo = false;
            }
            if new_line != "" {
                new_file += &(new_line + "\n");
            }

            index += 1;
        }

        new_file
    }
}

fn main() {
    start_timer();

    let start_time = Instant::now();
    let args: Vec<String> = env::args().collect();

    if args.len() < 4 { // Rust args include the program name itself
        let elapsed = start_time.elapsed().as_secs_f64();
        eprintln!(
            "{}{}{}: {}{}Argument error:{} Not enough arguments {}(Must have 3: file path, output path, syntactic mode){}",
            FG_BRIGHT_CYAN, elapsed, RESET, BOLD, FG_BRIGHT_RED, RESET, FG_BRIGHT_BLUE, RESET
        );
        exit(1);
    } else if args.len() > 4 {
        let elapsed = start_time.elapsed().as_secs_f64();
        eprintln!(
            "{}{}{}: {}Too many arguments{}{} {}(Should have 3: file path, output path, syntactic mode){}",
            FG_BRIGHT_CYAN, elapsed, RESET, BOLD, FG_BRIGHT_YELLOW, RESET, FG_BRIGHT_BLUE, RESET
        );
    }

    let path: &str = &args[1].clone(); // Clone the path to own it.
    let output: Option<&str> = if args.len() > 1 && args[1] == "none" {
        None
    } else if args.len() > 1{
        Some(&args[1])
    } else {
        None
    };
    let syntactic: bool = if args.len() > 1 && args[1].to_lowercase() == "true" {
        true
    } else {
        false
    };

    let elapsed = start_time.elapsed().as_secs_f64();
    println!(
        "{}{}{}: {}Running with parameters {}:{}{} \n{}path: {}, output path: {:?}, syntactic mode: {}{}",
        FG_BRIGHT_CYAN, elapsed, RESET, FG_BRIGHT_GREEN, RESET, FG_BRIGHT_BLUE, RESET, FG_BRIGHT_BLUE, path,
        match output {
            Some(s) => {}
            None => {
                "";
            }
        },
        syntactic, RESET
    );

    let app: Result<App, String> = App::new(path, output, syntactic);

    match app {
        Ok(mut app) => {
            app.run(); // Correct: call run() on the App instance
        }
        Err(err) => {
            eprintln!("Error creating App: {}", err);
        }
    }
}
