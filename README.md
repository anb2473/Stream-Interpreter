# Stream Programming Language

Stream is a new, flexible, and portable multi-paradigm programming language designed for rapid development and ease of use. Stream aims to simplify complex programming tasks with a clean, functional syntax and robust cross-platform capabilities.

## What is Stream?

* **Ease of Development:** Stream prioritizes a straightforward development experience.
* **Portability:** Designed for cross-platform compatibility.
* **Flexible Syntax:** Offers a syntax that adapts to various programming styles.
* **Functional Approach:** Stream promotes a functional programming style, aiming for clear and maintainable code.
* **Python-Powered:** Built using Python for rapid iteration and development.

## Stream Syntax Examples
```python
# This is a comment

let x: int = 0;   # Semi colons are fully optional
y = 'c'

if y:str == 'a':char {
  return 1
}
elif y:str == 'b':char {
  return 2
}
else {
  let check = if not x:int == 3:int;
  while check {
    x = x + 1
    check = if not x:int == 3:int
  }
}

def navigate_dict -> str (value: int) {
  values = {1: 'a', 2: 'b', 3: 'c'}
  new_value = values.value:int; return new_value
}

return navigate_dict(1)   # Returns 'a'
```

## Getting Started (Temporary Instructions)

**Note:** We are actively developing an installer to simplify this process.

1.  **Install Python:** Download and install Python from (python.org/downloads).
2.  **Locate Compiler and Interpreter:** Find the `compiler.py` and `interpreter.py` files within the cloned repository.

### Compilation

* The Stream compiler takes a `.stream` file as input and produces a compiled file.
* Open your terminal or command prompt.
* Run the compiler using:

    ```powershell
    python "[PATH_TO_COMPILER]" "[PATH_TO_YOUR_STREAM_FILE]"
    ```

* **Important:** Remember to use `"` marks, for example:

    ```powershell
    python "C:\My Projects\Stream\compiler.py" "C:\My Projects\Stream\my program.stream"
    ```

### Interpretation

* The Stream interpreter executes the compiled file.
* After compilation, run the interpreter using:

    ```powershell
    python "[PATH_TO_INTERPRETER]" "[PATH_TO_COMPILED_FILE]" "[RETURN_TYPE]" "{}"
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
