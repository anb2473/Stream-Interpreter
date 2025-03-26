# Stream Programming Language

Stream is a new, flexible, and portable multi-paradigm programming language designed for rapid development and ease of use. Stream aims to simplify complex programming tasks with a clean, functional syntax and robust cross-platform capabilities.

**PLEASE GO TO OUR WIKI FOR MORE INFORMATION ON LANGUAGE SYNTAX AND INSTALLATION FOR YOUR SPECIFIC VERSION, OR HEAD TO OUR WEBSITE (stream-compiler.netlify.app)**

## What is Stream?

* **Ease of Development:** Stream prioritizes a straightforward development experience.
* **Portability:** Designed for cross-platform compatibility.
* **Flexible Syntax:** Offers a syntax that adapts to various programming styles.
* **Functional Approach:** Stream promotes a functional programming style, aiming for clear and maintainable code.
* **Python-Powered:** Built using Python for rapid iteration and development.

## Stream Syntax Examples

* Check out our wiki for information on the language features and syntax.
* Please try to run our test program in the `TEST.md` file. If you do not get the same output, please reach out in the Issues section.

## Getting Started (Temporary Instructions)

**Note:** We are actively developing an installer to simplify this process.

### THE FOLLOWING INFORMATION IS FOR THE NEWEST MODEL OF STREAM (v1.2.0+)

Please go to our wiki for more detailed install information on specific releases.

1.  **Install Python:** Download and install Python from (python.org/downloads).
2.  **Locate Compiler and Interpreter:** Find the `compiler.py` and `interpreter.py` files within the cloned repository.

### Compilation

* The Stream compiler takes a `.stream` file as input and produces a compiled file with a `.fsl` extension.
* Open your terminal or command prompt.
* Run the compiler using:

    ```powershell
    python "[PATH_TO_COMPILER]" "[PATH_TO_YOUR_STREAM_FILE.stream]"
    ```

* **Important:** Remember to use `"` marks, for example:

    ```powershell
    python "C:\My Projects\Stream\compiler.py" "C:\My Projects\Stream\my program.stream"
    ```

### Interpretation

* The Stream interpreter executes the compiled file.
* After compilation, run the interpreter using:

    ```powershell
    python "[PATH_TO_INTERPRETER]" "[PATH_TO_COMPILED_FILE.fsl]" "[RETURN_TYPE]" "{}"
    ```

* **Important:** Remember to use `"` marks.
* **Note:** The `[RETURN_TYPE]` specifies the expected return type (e.g., `int`, `str`). The `{}` represents JSON parameters, which are typically empty for basic execution.

## Simplified Commands (Windows PowerShell ONLY)

To make this easier on Windows, you can add functions to your PowerShell profile:

1.  **Open your PowerShell profile:**
    * Run: `notepad $PROFILE`
    * If that doesn't work, run `New-Item -Path $PROFILE -Type File -Force` first.
2.  **Add these functions to your profile:**

    ```powershell
    function compile_stream {
        param(
            [Parameter(Mandatory=$true)]
            [string]$file
        )
        python "[PATH_TO_COMPILER]" $file
    }

    function interpret_stream {
        param(
            [Parameter(Mandatory=$true)]
            [string]$file,
            [string]$type,
            [string]$json
        )
        python "[PATH_TO_INTERPRETER]" $file $type $json
    }
    ```

    Now, you can use:

    * `compile_stream " [PATH_TO_STREAM_FILE]"`
    * `interpret_stream " [PATH_TO_COMPILED_FILE]" "[RETURN_TYPE]" "{}"`

## System Details

* **Intermediate Compilation:** Stream compiles to a type-safe intermediate file before final execution or compilation.
* **Python Implementation:** The compiler and interpreter are written in Python for rapid development.

## Next Steps

* **Cross-Platform Compilation:** Implement compilation to C, Java, Python, and Rust.
* **Stream IDE:** Develop a dedicated IDE for Stream.
* **Web Integration:** Add web capabilities to Stream.
* **Installers:** Create installers for easy setup.

## Reach Out

* If you wish to reach out to us, please feel free to contact us at our email, streamissues@gmail.com.
* If you wish to report an issue or bug, please enter it under the Issues section.
* For any major security concerns, please read our instructions in the SECURITY.md file. Please **DO NOT** post major security issues in the Issues section or send them publicly until the issue has been resolved.
  
## For More Details

Please read our wiki page for more details on language syntax, next steps, or any other issues.
