import os
import shutil
import time
from cfonts import render
from pynput import keyboard

red = "\033[91m"
orange = "\033[93m"
yellow = "\033[33m"
green = "\033[92m"
blue = "\033[94m"
indigo = "\033[96m"
violet = "\033[95m"
magenta = "\033[95m"

reset = "\033[0m"


class ArrowKeyHandler:
    def __init__(self):
        self.key_queue = []
        self.lock = threading.Lock()
        self.listener = None
        self.start_listener()

    def on_press(self, key):
        with self.lock:
            if key == keyboard.Key.up:
                self.key_queue.append("up")
            elif key == keyboard.Key.down:
                self.key_queue.append("down")
            elif key == keyboard.Key.enter:
                self.key_queue.append("enter")

    def start_listener(self):
        """Start the keyboard listener in a non-blocking way"""
        if self.listener is None or not self.listener.running:
            self.listener = keyboard.Listener(on_press=self.on_press)
            self.listener.daemon = True
            self.listener.start()

    def get_key(self):
        """Get the next key from the queue"""
        with self.lock:
            if self.key_queue:
                return self.key_queue.pop(0)
            return None

    def stop_listener(self):
        """Stop the keyboard listener"""
        if self.listener is not None and self.listener.running:
            self.listener.stop()
            self.listener = None


# Global handler that can be reused across multiple calls
_arrow_key_handler = None


def arrow_key():
    """
    Detects arrow key presses and returns:
    - True for up arrow
    - False for down arrow
    - 0 for enter
    - None if no key was pressed
    """
    global _arrow_key_handler

    # Initialize the handler if it doesn't exist
    if _arrow_key_handler is None:
        _arrow_key_handler = ArrowKeyHandler()

    # Process a small timeout to allow for key detection
    # without completely blocking the application
    import time
    time.sleep(0.05)

    # Get the latest key press
    key = _arrow_key_handler.get_key()

    if key == "up":
        return True
    elif key == "down":
        return False
    elif key == "enter":
        return 0
    else:
        return None


# For proper cleanup when your application is done
def cleanup_arrow_key_handler():
    global _arrow_key_handler
    if _arrow_key_handler is not None:
        _arrow_key_handler.stop_listener()
        _arrow_key_handler = None


def cls():
    os.system('cls' if os.name == 'nt' else 'clear')


def render_choice(question, options, index):
    print(question)
    colors = [red, orange, yellow, green, blue, indigo, violet, magenta]
    for i in range(len(options)):
        if i == index:
            print(f"{colors[i % 4 + 1]}>  {options[i]}{reset}")
        else:
            print(f"{colors[i % 4]}   {options[i]}{reset}")


def render_question(question, options):
    index = 0
    cls()
    render_choice(question, options, index)
    arrow_key()
    time.sleep(0.25)
    while True:
        ret = arrow_key()
        if ret is not None:
            if ret is True:
                index = max(index - 1, 0)
            elif ret is False:
                index = min(index + 1, len(options) - 1)
            else:
                return index
            cls()
            render_choice(question, options, index)


def get_files(directory_path='.'):
    directories = []
    files = []

    try:
        # Use scandir() instead of listdir() for better performance and metadata access
        with os.scandir(directory_path) as entries:
            for entry in entries:
                try:
                    # is_dir() is more reliable than os.path.isdir() as it uses file descriptor info
                    if entry.is_dir():
                        directories.append(entry.name)
                    else:
                        files.append(entry.name)
                except OSError:
                    # Handle case where file might be deleted during iteration
                    print(f"{yellow}Warning: Could not access {entry.name}{reset}")
    except FileNotFoundError:
        print(f"{red}Error: Directory '{directory_path}' does not exist.{reset}")
    except PermissionError:
        print(f"{red}Error: No permission to access '{directory_path}'.{reset}")
    except NotADirectoryError:
        print(f"{red}Error: '{directory_path}' is not a directory.{reset}")

    return directories, files


def read_file(filepath):
    try:
        # Attempt to open the file in read mode ('r')
        with open(filepath, 'r') as file:
            # Read the entire content of the file
            file_content = file.read()
            return file_content  # Return the content if successful

        # Handle specific exceptions that might occur during file operations
    except FileNotFoundError:
        print(f"Error: File not found at {filepath}")
        return None  # Explicitly return None to indicate failure
    except PermissionError:
        print(f"Error: Permission denied to read {filepath}")
        return None  # Explicitly return None to indicate failure
    except IsADirectoryError:
        print(f"Error: {filepath} is a directory, not a file")
        return None
    except UnicodeDecodeError:
        print(f"Error: {filepath} is not a text file encoded in a supported format (e.g., UTF-8)")
        return None
    except OSError as e:
        print(f"Error reading {filepath}: {e}")
        return None  # Catch other OS-related errors
    except Exception as e:
        # Catch any other unexpected exceptions
        print(f"An unexpected error occurred while reading {filepath}: {e}")
        return None  # Ensure None is returned on any error


def get_console_height():
    try:
        # Use shutil.get_terminal_size() which is more robust
        terminal_size = shutil.get_terminal_size()
        return terminal_size.lines
    except OSError:
        print("Error: Could not get console size. Defaulting to 24.")
        return 24  # Default height if the terminal size cannot be determined
    except Exception as e:
        print(f"An unexpected error occurred while getting console height: {e}")
        return 24


import os
from pynput import keyboard
import threading


def get_console_height():
    """Get terminal height"""
    try:
        return os.get_terminal_size().lines
    except:
        return 24  # Fallback height


def read_file(file_path):
    """Read file content"""
    try:
        with open(file_path, 'r') as f:
            return f.read()
    except:
        return ""


class KeyHandler:
    def __init__(self):
        self.key_queue = []
        self.lock = threading.Lock()
        self.listener = None
        self.start_listener()

    def on_press(self, key):
        with self.lock:
            if key == keyboard.Key.up:
                self.key_queue.append("up")
            elif key == keyboard.Key.down:
                self.key_queue.append("down")
            elif key == keyboard.Key.left:
                self.key_queue.append("left")
            elif key == keyboard.Key.right:
                self.key_queue.append("right")
            elif key == keyboard.Key.enter:
                self.key_queue.append("enter")
            elif key == keyboard.Key.backspace:
                self.key_queue.append("backspace")
            elif key == keyboard.Key.delete:
                self.key_queue.append("delete")
            elif key == keyboard.Key.home:
                self.key_queue.append("home")
            elif key == keyboard.Key.end:
                self.key_queue.append("end")
            elif key == keyboard.Key.esc:
                self.key_queue.append("escape")
            elif key == keyboard.Key.tab:
                self.key_queue.append("tab")
            elif hasattr(key, 'char'):
                if key.char == '\x13':  # Ctrl+S
                    self.key_queue.append("ctrl_s")
                elif key.char == '\x11':  # Ctrl+Q
                    self.key_queue.append("ctrl_q")
                elif key.char is not None:
                    self.key_queue.append(key.char)

    def start_listener(self):
        """Start the keyboard listener in a non-blocking way"""
        if self.listener is None or not self.listener.running:
            self.listener = keyboard.Listener(on_press=self.on_press)
            self.listener.daemon = True
            self.listener.start()

    def get_key(self):
        """Get the next key from the queue"""
        with self.lock:
            if self.key_queue:
                return self.key_queue.pop(0)
            return None

    def stop_listener(self):
        """Stop the keyboard listener"""
        if self.listener is not None and self.listener.running:
            self.listener.stop()
            self.listener = None


def kit(file_path):
    # Check if file exists, if not create it
    if not os.path.exists(file_path):
        with open(file_path, 'w') as f:
            f.write('')

    # Read the file content
    content = read_file(file_path)
    lines = content.split('\n')

    # Initialize variables
    current_line = 0
    offset = 0
    edit_mode = False
    edit_line = ''
    edit_position = 0
    message = f"Viewing: {file_path} | Press Enter to edit line | Ctrl+S to save | Ctrl+Q to quit"

    # Initialize key handler
    key_handler = KeyHandler()

    try:
        while True:
            # Clear screen
            os.system('cls' if os.name == 'nt' else 'clear')

            # Get terminal height
            terminal_height = get_console_height() - 2  # Reserve space for message and prompt

            # Calculate visible range
            max_offset = max(0, len(lines) - terminal_height)
            if offset > max_offset:
                offset = max_offset

            end_index = min(offset + terminal_height, len(lines))

            # Display file content
            for i in range(offset, end_index):
                prefix = '> ' if i == current_line and not edit_mode else '  '
                if i < len(lines):
                    print(f"{prefix}{i + 1}: {lines[i]}")
                else:
                    print(f"{prefix}{i + 1}:")

            # Display message
            print(message)

            if edit_mode:
                # Display edit prompt
                edit_display = edit_line[:edit_position] + '|' + edit_line[edit_position:]
                print(f"Edit line {current_line + 1}: {edit_display}")

            # Process any key presses in the queue
            key = None
            # Wait a bit for key input to accumulate and improve responsiveness
            import time
            time.sleep(0.05)  # Small sleep to prevent CPU overload

            # Get latest key press
            while True:
                next_key = key_handler.get_key()
                if next_key is None:
                    break
                key = next_key  # Use the most recent key press

            if key is None:
                continue

            if edit_mode:
                # Handle edit input
                if key == "ctrl_s":  # Save
                    lines[current_line] = edit_line
                    save_file(file_path, lines)
                    edit_mode = False
                    message = f"Line {current_line + 1} saved. Viewing: {file_path}"
                elif key == "ctrl_q" or key == "escape":  # Quit edit mode
                    edit_mode = False
                    message = f"Edit canceled. Viewing: {file_path}"
                elif key == "backspace":  # Backspace
                    if edit_position > 0:
                        edit_line = edit_line[:edit_position - 1] + edit_line[edit_position:]
                        edit_position -= 1
                elif key == "delete":  # Delete
                    if edit_position < len(edit_line):
                        edit_line = edit_line[:edit_position] + edit_line[edit_position + 1:]
                elif key == "left":  # Left arrow
                    if edit_position > 0:
                        edit_position -= 1
                elif key == "right":  # Right arrow
                    if edit_position < len(edit_line):
                        edit_position += 1
                elif key == "home":  # Home
                    edit_position = 0
                elif key == "end":  # End
                    edit_position = len(edit_line)
                elif key == "enter":  # Enter - confirm edit
                    lines[current_line] = edit_line
                    edit_mode = False
                    message = f"Line {current_line + 1} edited. Press Ctrl+S to save changes."
                elif key is not None and len(key) == 1:  # Regular character
                    edit_line = edit_line[:edit_position] + key + edit_line[edit_position:]
                    edit_position += 1
            else:
                # Handle navigation input
                if key == "up":  # Up arrow
                    if current_line > 0:
                        current_line -= 1
                        if current_line < offset:
                            offset = current_line
                elif key == "down":  # Down arrow
                    if current_line < len(lines) - 1:
                        current_line += 1
                        if current_line >= offset + terminal_height:
                            offset += 1
                elif key == "enter":  # Enter key
                    edit_mode = True
                    edit_line = lines[current_line] if current_line < len(lines) else ""
                    edit_position = len(edit_line)
                    message = f"Editing line {current_line + 1} | Ctrl+S to save | Ctrl+Q/ESC to cancel"
                elif key == "ctrl_q":  # Quit
                    break
                elif key == "ctrl_s":  # Save
                    save_file(file_path, lines)
                    message = f"File saved: {file_path}"
    finally:
        # Always clean up the keyboard listener
        key_handler.stop_listener()


def save_file(file_path, lines):
    """Save the content to the file"""
    with open(file_path, 'w') as f:
        f.write('\n'.join(lines))




class CLI:
    def __init__(self):
        cls()

        self.root = "C:\\"
        self.cmd = ""

        print(f"""
    {red} ________   _________   ________   _______    ________   _____ ______{reset}
    {orange}|\\   ____\\ |\\___   ___\\|\\   __  \\ |\\  ___ \\  |\\   __  \\ |\\   _ \\  _   \\{reset}
    {yellow}\\ \\  \\___|_\\|___ \\  \\_|\\ \\  \\|\\  \\\\ \\   __/| \\ \\  \\|\\  \\\\ \\  \\\\\\__\\ \\  \\{reset}
    {green} \\ \\_____  \\    \\ \\  \\  \\ \\   _  _\\\\ \\  \\_|/__\\ \\   __  \\\\ \\  \\\\|__| \\  \\{reset}
    {blue}  \\|____|\\  \\    \\ \\  \\  \\ \\  \\\\  \\|\\ \\  \\_|\\ \\\\ \\  \\ \\  \\\\ \\  \\    \\ \\  \\{reset}
    {indigo}    ____\\_\\  \\    \\ \\__\\  \\ \\__\\\\ _\\ \\ \\_______\\\\ \\__\\ \\__\\\\ \\__\\    \\ \\__\\{reset}
    {violet}   |\\_________\\    \\|__|   \\|__|\\|__| \\|_______| \\|__|\\|__| \\|__|     \\|__|{reset}
    {magenta}   \\|_________|{reset}        __            __    _  __ 
                       ___/ /___  _  __ / /__ (_)/ /_
                      / _  // -_)| |/ //  '_// // __/
                      \\_,_/ \\__/ |___//_/\\_\\/_/ \\__/ 
    """)

        print("        Copyright © 2025 Austin Nabil Blass. All rights reserved.\n")

        time.sleep(.5)

        print("     Welcome to the Stream devkit. This is a CLI interface for working with\n"
              "Stream tools. For help, please type 'help.'\n")

    def rm(self, item):
        item = f"{self.root}\\{item}"
        if not os.path.isdir(item):
            try:
                os.remove(item)
                print(f"{green}Successfully removed file: {item}{reset}")
            except FileNotFoundError:
                print(f"{yellow}File not found at: {item}{reset}")
            except PermissionError:
                print(f"{red}Permission denied to remove file: {item}{reset}")
            except IsADirectoryError:
                print(
                    f"{red}'{item}' is a directory, failed to remove{reset}")
            except OSError as e:
                print(f"{red}Error removing file {item}: {e}{reset}")
            except Exception as e:
                print(f"{red}An unexpected error occurred while trying to remove {item}: {e}{reset}")
        else:
            try:
                shutil.rmtree(item)
                print(f"{green}Successfully removed directory: {item}{reset}")
            except FileNotFoundError:
                print(f"{yellow}Directory not found at: {item}{reset}")
            except OSError as e:
                print(f"{red}Error removing directory {item}: {e}{reset}")

    def mk(self, item):
        if not item.startswith("C:\\"):
            item = f"{self.root}\\{item}"
        try:
            if not item.__contains__('.'):
                if os.path.exists(item):
                    print(f"{yellow}Directory already exists: {item}{reset}")
                    return

                os.makedirs(item)
                print(f"{green}Successfully created directory: {item}{reset}")
            else:
                if os.path.exists(item):
                    print(f"{yellow}File already exists: {item}{reset}")
                    return

                dir_path = os.path.dirname(item)
                if dir_path and not os.path.exists(dir_path):
                    os.makedirs(dir_path)

                with open(item, 'x'):
                    print(f"{green}Successfully created file: {item}{reset}")

        except FileExistsError:
            print(f"{yellow}File or directory already exists at: {item}{reset}")
        except PermissionError:
            print(f"{red}Permission denied to create item at: {item}{reset}")
        except OSError as e:
            print(f"{red}Error creating item at {item}: {e}{reset}")
        except Exception as e:
            print(f"{red}An unexpected error occurred while creating item at {item}: {e}{reset}")

    def help(self):
        ret = render_question("Please select what you need help with",
                              ["STREAM DEVKIT", "CLI TOOLS", "STREAM SYNTAX", "EXIT"])
        if ret == 3:
            return
        if ret == 0:
            print("\nSTREAM DEVKIT\n"
                  "The Stream devkit allows you to access many builtin Stream and CLI tools via a CLI interface.\n"
                  "Type 'new <your_app_name>' to create a new project. Run 'compile <your_stream_file> <output_path> \n"
                  "<syntactic_mode>' to compile your stream file, and 'exec <your_stream_file> <return_type> \n"
                  "<params_json> <verbose_mode>' to execute the generated stream file. (replace output_path with none\n"
                  "to automatically generate a path) Use 'build <your_stream_file>' to compile and execute your stream\n"
                  "file in one command.\n\n")
            return
        elif ret == 1:
            print(
                "CLI TOOLS\n"
                "For basic CLI tools, type 'cd <dir_name>' to open a directory. Replace dir_name with '..' to \n"
                "exit a directory. Type a 'C:\\' path to go directly to that path. Run 'ls' to list the files \n"
                "in the current directory. Run 'mk <file_name>' to make a file, and 'rm <file_name>' to remove\n"
                "a file. (rm and mk work for both directories and files) Type 'kit <file_name>' \n"
                "to automatically open the devkit file editor. Type 'exit' to close the devkit. Type 'clear'\n"
                "to clear the terminal.\n")
            return
        else:
            ret = render_question("Please select what you want to learn",
                                  ["OVERVIEW", "VARIABLES", "CONDITIONALS & LOOPS", "LOOPS", "FUNCTIONS", "BACK"])
            if ret == 5:
                self.help()
            answers = [
                "OVERVIEW\n"
                "The Stream language syntax is designed to be modern and concise.",
                "VARIABLES\n"
                "The Stream language provides many ways of declaring variables. You can use 'let x: int = 0' to declare\n"
                "a variable. However, Stream provides type less alternatives, allowing you to instead simply write\n"
                "'x = 0'. This type inference works for declaring raw variables like you just saw, declaring action\n"
                "variables such as 'x = y + 1', and declaring function variables such as 'x = foo()'. You can also use\n"
                "the 'obj' type to use Python builtin variables. For example, if you want to call a python function\n"
                "which returns an unknown type, use the 'obj' type and the interpreter will automaticaly register the\n"
                "type.\n",
                "CONDITIONALS\n"
                "The Stream language is designed to be procedural / functional, and as such offers both procedural and\n"
                "traditional methods of declaring conditionals. You can simply type 'if x <= 0 { ... }' to write\n"
                "a simple conditional. You can also use 'elif' and 'else' for more complex logic. These elements\n"
                "will automatically compile down to a procedural format using 'do' blocks. These blocks are written\n"
                "as 'do True 3', where the interpreter will check if the first parameter is False, and if so jump\n"
                "the next n lines, declared by the second variable. (n can also be negative to allow more complex\n"
                "logic) You can also write conditionals as variables. For example, 'x = if n == True'.\n\n"
                "LOOPS\n"
                "Loops can be similarly written in two ways. You can write 'while check = if x <= 0 { ... }'.\n"
                "You can also write them with 'while check { ... }' if check has been previously declared.\n"
                "And, you can also use 'do' blocks to represent loops, allowing you to target a more procedural\n"
                "paradigm.\n",
                "FUNCTIONS\n"
                "The language allows you to declare functions in three different ways. You can write 'prototype \n"
                "my_foo C:/path/to/my_foo' -> int (x: int)'. This is useful for importing functions from other\n"
                "projects and also allows you to easily import python functions and functions from other\n"
                "languages. You can also simply write 'def my_foo -> int (x: int)' to declare a function.\n"
                "This will similarly be compiled to a prototype statement. However, if your goal is to create\n"
                "a one file program, you can instead use 'do' blocks to jump between code segments, however\n"
                "this is not the recommended approach.\n"
            ]

            print(answers[ret])

    def exec(self):
        func = self.cmd.strip().split(' ', 1)[0]
        match func:
            case "cd":
                o_root = self.root
                self.cmd = self.cmd[3:].strip()
                if self.cmd.startswith("C:\\"):
                    self.root = self.cmd
                elif self.cmd.startswith(".."):
                    if self.root != "C:\\":
                        self.root = self.root.rsplit("\\", 1)[0]
                else:
                    self.root += f"\\{self.cmd}" if not self.root.endswith("\\") else self.cmd
                if not os.path.exists(self.root):
                    self.root = o_root
                    print(f"{red}Cannot open directory, path not found{reset}")
            case "compile":
                param = self.cmd.split(" ")[1:]
                os.system(f"./C:\\Stream\\bin\\stream-fsl.exe "
                          f"{param[0] if param[0].startswith("C:\\") else f"{self.root}\\{param[0]}"} "
                          f"{param[1] if param[1].startswith("C:\\") else f"{self.root}\\{param[1]}"} "
                                                                          f"{param[2]}")
            case "exec":
                param = self.cmd.split(" ")[1:]
                os.system(f"./C:\\Stream\\bin\\stream-e.exe "
                          f"{param[0] if param[0].startswith("C:\\") else f"{self.root}\\{param[0]}"} "
                          f"{param[1]} {param[2]} {param[3]}")
            case "build":
                param = self.cmd.split(" ")[1:]
                os.system(f"./C:\\Stream\\bin\\stream-fsl.exe "
                          f"{param[0] if param[0].startswith("C:\\") else f"{self.root}\\{param[0]}"} "
                          f"{param[1] if param[1].startswith("C:\\") or param[1] == "none" else f"{self.root}\\{param[1]}"} "
                          f"{param[2]}")
                path = param[1] if not param[1] == "none" else (f"{param[0].split(".")[0] if param[0].startswith("C:\\") else self.root}\\build\\build_"
                                                                f"{param[0].rsplit("\\", 1)[len(param[0].rsplit("\\", 1)) - 1]}"
                                                                f"\\{param[0].rsplit("\\", 1)[len(param[0].rsplit("\\", 1)) - 1]}.fsl")
                os.system(f"./C:\\Stream\\bin\\stream-e.exe "
                          f"{path} {param[3]} {param[4]} {param[5]}")
            case "title":
                print(f"""
                        {red} ________   _________   ________   _______    ________   _____ ______{reset}
                        {orange}|\\   ____\\ |\\___   ___\\|\\   __  \\ |\\  ___ \\  |\\   __  \\ |\\   _ \\  _   \\{reset}
                        {yellow}\\ \\  \\___|_\\|___ \\  \\_|\\ \\  \\|\\  \\\\ \\   __/| \\ \\  \\|\\  \\\\ \\  \\\\\\__\\ \\  \\{reset}
                        {green} \\ \\_____  \\    \\ \\  \\  \\ \\   _  _\\\\ \\  \\_|/__\\ \\   __  \\\\ \\  \\\\|__| \\  \\{reset}
                        {blue}  \\|____|\\  \\    \\ \\  \\  \\ \\  \\\\  \\|\\ \\  \\_|\\ \\\\ \\  \\ \\  \\\\ \\  \\    \\ \\  \\{reset}
                        {indigo}    ____\\_\\  \\    \\ \\__\\  \\ \\__\\\\ _\\ \\ \\_______\\\\ \\__\\ \\__\\\\ \\__\\    \\ \\__\\{reset}
                        {violet}   |\\_________\\    \\|__|   \\|__|\\|__| \\|_______| \\|__|\\|__| \\|__|     \\|__|{reset}
                        {magenta}   \\|_________|{reset}        __            __    _  __ 
                                           ___/ /___  _  __ / /__ (_)/ /_
                                          / _  // -_)| |/ //  '_// // __/
                                          \\_,_/ \\__/ |___//_/\\_\\/_/ \\__/ 
                        """)
            case "clear":
                cls()
            case "ls":
                files = get_files(self.root)
                print(f"{magenta}Sub Directories: [{reset}")
                index = 0
                length = 0
                for file in files[0]:
                    if index % 3 != 0:
                        if index % 3 != 2:
                            print(" " * max(30 - length, 3) + blue + file + reset, end="")
                        else:
                            print(" " * max(30 - length, 3) + blue + file + reset)
                        length = min(30 - length, 3) + len(file)
                    else:
                        length = len(file) + 5
                        print(f"     {blue + file + reset}", end="")
                    index += 1
                print(("" if index % 3 == 0 else "\n") + f"{magenta}]{reset}")
                print(f"{magenta}Files: [{reset}")
                index = 0
                for file in files[1]:
                    if index % 3 != 0:
                        if index % 3 != 2:
                            print(" " * max(30 - length, 3) + green + file + reset, end="")
                        else:
                            print(" " * max(30 - length, 3) + green + file + reset)
                        length = min(30 - length, 3) + len(file)
                    else:
                        length = len(file) + 5
                        print(f"     {green + file + reset}", end="")
                    index += 1
                print(("" if index % 3 == 0 else "\n") + f"{magenta}]{reset}")
            case "kit":
                print(f"""
                                      {yellow}|{reset}
                   {blue}.__{reset}                {orange}_{reset}		
                   {blue}|X |{reset}           {yellow}-{reset}  {red}( ){reset}  {yellow}-{reset}	
                   {blue}`-.|{reset}               {orange}"{reset}
                  '    {blue}`+_ +{reset}          {yellow}|{reset}
                .'           {indigo}`x{reset}
              .'                {indigo}={reset}
     {green}O{reset}    {magenta}_.-'     ..{reset}
    {green}( >-='{reset}            {magenta}``.._                  _.{reset}
{magenta}-'{reset}{green}  / \\"{reset}                   {magenta}`-a:f,-'-._    ,-'{reset}
    {green}()){reset}
  {magenta}_{reset} {green}b b{reset}      ____  __..__   __   
   {magenta}`"{reset}       |    |/ _||__|_/  |_   
            |      <  |  |\\   __\\\\
            |    |  \\ |  | |  |  
            |____|__ \\|__| |__|  
                    \\/ 
                """)

                time.sleep(.5)

                try:
                    kit(self.cmd.split(" ")[1].strip())
                except IndexError:
                    print(f"{red}Please specify a path{reset}")
            case "rm":
                self.rm(self.cmd[3:].strip())
            case "mk":
                self.mk(self.cmd[3:].strip())
            case "new":
                project_name = self.cmd[3:].strip()
                if project_name.__contains__('.'):
                    print(f"{red}Project name is a filepath, not a directory{reset}")
                else:
                    self.mk(project_name)
                    self.mk(project_name + "\\src")
                    with open(project_name + "\\src\\main.stream") as f:
                        f.write("return 0")
                    with open(project_name + "\\function_manager.spk") as f:
                        f.write("src\\main.fsl <%$> int ()")
                    self.root += project_name + "\\src"
                    print(f"{green}Project successfully created{reset}")
            case "help":
                self.help()
            case "exit":
                raise Exception()
            case default:
                if func != '':
                    print(f"{red}Unknown command{reset}")

        self.cmd = ""

    def get_input(self):
        self.cmd = input()
        if self.cmd.strip() == "":
            self.get_input()

    def run(self):
        try:
            print(f"{blue}> {violet}{{ {self.root} }}{reset} {green}>{reset} ", end="")
            self.get_input()
            self.exec()
            self.run()
        except Exception as e:
            print(e)
            print(f"""
        {blue}__{reset}{red}/\\\\\\\\\\\\\\\\\\\\\\\\\\{reset}{blue}____{reset}{red}/\\\\\\{reset}{blue}________{reset}{red}/\\\\\\{reset}{blue}__{reset}{red}/\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\{reset}{blue}_{reset}
         {blue}_{reset}{red}\\/\\\\\\/////////\\\\\\{reset}{blue}_{reset}{red}\\///\\\\\\{reset}{blue}____{reset}{red}/\\\\\\/{blue}__{reset}{red}\\/\\\\\\///////////{reset}{reset}{blue}__{reset}
          {blue}_{reset}{red}\\/\\\\\\{reset}{blue}_______{reset}{red}\\/\\\\\\{reset}{blue}___{reset}{red}\\///\\\\\\/\\\\\\/{blue}____{reset}{red}\\/\\\\\\{reset}{blue}_____________{reset}
           {blue}_{reset}{red}\\/\\\\\\\\\\\\\\\\\\\\\\\\\\\\{reset}{blue}______{reset}{red}\\///\\\\\\/{blue}______{reset}{red}\\/\\\\\\\\\\\\\\\\\\\\\\{reset}{blue}_____{reset}
            {blue}_{reset}{red}\\/\\\\\\/////////\\\\\\{reset}{blue}_______{reset}{red}\\/\\\\\\{blue}_______{reset}{red}\\/\\\\\\///////{reset}{blue}______{reset}
             {blue}_{reset}{red}\\/\\\\\\{reset}{blue}_______{reset}{red}\\/\\\\\\{reset}{blue}_______{reset}{red}\\/\\\\\\{blue}_______{reset}{red}\\/\\\\\\{reset}{blue}_____________{reset}
              {blue}_{reset}{red}\\/\\\\\\{reset}{blue}_______{reset}{red}\\/\\\\\\{reset}{blue}_______{reset}{red}\\/\\\\\\{blue}_______{reset}{red}\\/\\\\\\{reset}{blue}_____________{reset}
               {blue}_{reset}{red}\\/\\\\\\\\\\\\\\\\\\\\\\\\\\/{blue}________{reset}{red}\\/\\\\\\{blue}_______{reset}{red}\\/\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\{reset}{blue}_{reset}
                {blue}_{reset}{red}\\/////////////{reset}{blue}__________{reset}{red}\\///{reset}{blue}________{reset}{red}\\///////////////{reset}{blue}__{reset}
""")


if __name__ == '__main__':
    app = CLI()
    app.run()
