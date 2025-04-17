import os
import sys

# Define print variables
RESET = "\033[0m"
BOLD = "\033[1m"
ITALIC = "\033[3m"

FG_BRIGHT_RED = "\033[91m"
FG_BRIGHT_GREEN = "\033[92m"
FG_BRIGHT_YELLOW = "\033[93m"
FG_BRIGHT_BLUE = "\033[94m"
FG_BRIGHT_CYAN = "\033[96m"


def compile():
    if len(sys.argv) < 2:
        print(f'{BOLD}{FG_BRIGHT_RED}Argument error:{RESET} '
              f'Not enough arguments {FG_BRIGHT_BLUE}(Must have 1: file path){RESET}')
        sys.exit(1)
    elif len(sys.argv) > 2:
        print(f'{BOLD}{FG_BRIGHT_YELLOW}'
              f'Too many arguments{RESET} {FG_BRIGHT_BLUE}(Should have 1: file path){RESET}')

    f = str(sys.argv[1])

    with open(f, 'r', encoding='utf-8') as file:
        content = file.read()

        new_content = ""

        split_content = content.split("{")

        for section in split_content:
            if not section.__contains__("}"):
                new_content += section + "{"
            else:
                value = section.split("}")[0]
                etc = section.split("}")[1]

                ret = value.split("]")[0].split("[")[1]

                value = value.split("]")[1]

                with open("temp.stream", 'w') as fl:
                    fl.write(value)

                os.system(f"stream-c temp.stream none false")

                with open(os.path.abspath("temp.stream"), 'r') as fl:
                    new_value = f"\n[{ret}]{fl.read()}"

                new_content += str(new_value) + "}" + etc + "{"

    with open(f.split(".")[0] + ".html", 'w') as file:
        file.write(new_content)


if __name__ == '__main__':
    compile()
