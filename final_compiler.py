import os
import re
import sys
import threading
import time

RESET = "\033[0m"       # Resets all formatting to default
BOLD = "\033[1m"        # Makes text bold
DIM = "\033[2m"         # Makes text dim/faint
ITALIC = "\033[3m"      # Makes text italic (not widely supported)
UNDERLINE = "\033[4m"   # Underlines text
BLINK = "\033[5m"       # Makes text blink (not widely supported)
REVERSE = "\033[7m"     # Reverses foreground and background colors
HIDDEN = "\033[8m"      # Makes text hidden/invisible

# Foreground Colors
FG_BLACK = "\033[30m"
FG_RED = "\033[31m"
FG_GREEN = "\033[32m"
FG_YELLOW = "\033[33m"
FG_BLUE = "\033[34m"
FG_MAGENTA = "\033[35m"
FG_CYAN = "\033[36m"
FG_WHITE = "\033[37m"

# Bright Foreground Colors
FG_BRIGHT_BLACK = "\033[90m"
FG_BRIGHT_RED = "\033[91m"
FG_BRIGHT_GREEN = "\033[92m"
FG_BRIGHT_YELLOW = "\033[93m"
FG_BRIGHT_BLUE = "\033[94m"
FG_BRIGHT_MAGENTA = "\033[95m"
FG_BRIGHT_CYAN = "\033[96m"
FG_BRIGHT_WHITE = "\033[97m"

# Background Colors
BG_BLACK = "\033[40m"
BG_RED = "\033[41m"
BG_GREEN = "\033[42m"
BG_YELLOW = "\033[43m"
BG_BLUE = "\033[44m"
BG_MAGENTA = "\033[45m"
BG_CYAN = "\033[46m"
BG_WHITE = "\033[47m"

# Bright Background Colors
BG_BRIGHT_BLACK = "\033[100m"
BG_BRIGHT_RED = "\033[101m"
BG_BRIGHT_GREEN = "\033[102m"
BG_BRIGHT_YELLOW = "\033[103m"
BG_BRIGHT_BLUE = "\033[104m"
BG_BRIGHT_MAGENTA = "\033[105m"
BG_BRIGHT_CYAN = "\033[106m"
BG_BRIGHT_WHITE = "\033[107m"

start_time = time.time()


# TODO: ADD MORE AUTOMATED TYPE GENERATION
# TODO: INTEGRATE SYSTEMS FOR CALLING PYTHON FILES AND MORE LANGUAGES
# TODO: ADD MORE FEATURES FOR CONDITIONAL STATEMENTS (<= between str and int)
# TODO: ADD MORE FEATURES FOR EXPRESSIONS
# TODO: INTEGRATE ERROR HANDLING
# TODO: IMPLEMENT GENERICS
# TODO: BUILD SEMANTIC ANALYZER / INTEGRATE SEMANTIC ANALYZER
# TODO: MAKE COMPILER MORE FLEXIBLE WITH SEMANTIC DIFFERENCES
# TODO: MAKE CODE MORE EFFICIENT
# TODO: BUILD IDE
# TODO: ALLOW WRITING IF STATEMENTS IN WHILE LOOPS
# TODO: INCREASE OPTIONS FOR LIST / DICT MANIPULATION
# TODO: ADD COMMENTS
# TODO: PACKAGE MANAGER BUILDER
# TODO: TEST COMPILER WITH TASKS SUCH AS FIZBUZZ
# TODO: BUILD A WAY TO RUN PROGRAMS THROUGH A PACKAGE MANAGER INSTEAD OF WITH LOTS OF COMMANDS
# TODO: CODE COMPILER IN STREAM


def read_file(path):
    try:
        # ATTEMPT TO READ FILE CONTENTS
        with open(path, 'r') as f:
            content = f.read()
            print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {FG_BRIGHT_GREEN}Successfully read file{RESET}'
                  f' at \'{path}\': \n{FG_BRIGHT_BLUE}{content}{RESET}')
            return content
    except FileNotFoundError:
        print(f'{FG_BRIGHT_CYAN}{time.time()}{RESET}: {FG_BRIGHT_RED}File Not Found:{RESET} \'{path}\'')
        return None
    except PermissionError:
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Permission denied:{RESET} '
              f'Unable to read file at \'{path}\'')
        return None
    except Exception as e:
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Unknown error:{RESET} '
              f'Unable to read file at \'{path}\': \'{e}\'')
        return None


def write_file(path, data):
    try:
        # ATTEMPT TO OVERRIDE DATA AT PATH / GENERATE NEW FILE
        with open(path, 'w') as f:
            f.write(data)
            print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {FG_BRIGHT_GREEN}Successfully wrote to file'
                  f'{RESET} at \'{path}\': \n{FG_BRIGHT_BLUE}{data}{RESET}')
            return 0
    except PermissionError:
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Permission denied:{RESET} '
              f'Unable to write to file at \'{path}\'')
        return 1
    except Exception as e:
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Unknown error:{RESET} '
              f'Unable to write to file at \'{path}\': \'{e}\'')
        return 2


def mkdir(path):
    try:
        # ATTEMPT TO GENERATE DIRECTORY
        os.mkdir(path)
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {FG_BRIGHT_GREEN}Successfully built directory{RESET}'
              f' at {ITALIC}\'{path}\'{RESET}')
        return 0
    except FileExistsError:
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}Directory already exists:'
              f'{RESET} {ITALIC}\'{path}\'{RESET}')
        return 0
    except PermissionError:
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Permission denied:{RESET}'
              f' Unable to build directory at {ITALIC}\'{path}\'{RESET}')
        return 2
    except Exception as e:
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Unknown error:{RESET}'
              f' Unable to build directory at {ITALIC}\'{path}\'{RESET}: {ITALIC}{FG_BRIGHT_RED}\'{e}\'{RESET}')
        return 3


def get_active_dir(path):
    return os.path.dirname(os.path.abspath(path))


class App(threading.Thread):
    def __init__(self, path, output=None, syntactic: bool = True):
        super().__init__()
        self.path: str = path
        self.output: str = output

        self.syntactic: bool = syntactic

        _, self.ext = os.path.splitext(self.path)

        # CHECK FOR VALID APPLICATION EXTENSION
        if self.ext != '.stream' and self.ext != '.sl':
            print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Invalid file extension:'
                  f'{RESET} '
                  f'{ITALIC}\'{self.ext}\'{RESET}')
            exit(1)

        # CHECK FILE PATH EXISTS
        if not os.path.exists(self.path):
            print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}File not found:{RESET} '
                  f'{ITALIC}\'{self.path}\'{RESET}')
            exit(1)

        # LOAD FILE CONTENTS
        self.file: str = read_file(self.path)

        # SAFELY CLOSE APPLICATION IF LOADING FILE FAILED
        if self.file is None:
            exit(1)

        # FIND WORKING DIRECTORY OF THE PATH
        self.active_directory: str = get_active_dir(path)

        # BUILD 'build' DIRECTORY
        ret: int = mkdir(f'{self.active_directory}\\build')

        # SAFELY CLOSE APPLICATION IF BUILD FAILED
        if ret != 0:
            exit(ret)

        # IF NO OUTPUT DIRECTORY IS SPECIFIED, GENERATE AN OUTPUT DIRECTORY
        if self.output is None:
            self.output = (f'{self.active_directory}\\build\\build_{os.path.basename(self.path).split('.')[0]}\\'
                           f'{os.path.basename(self.path).split('.')[0]}')

        # _, ext = os.path.splitext(self.path)
        # if ext == '.stream'

        # BUILD SUB BUILD DIRECTORY
        ret: int = mkdir(f'{self.active_directory}\\build\\build_{os.path.basename(self.path).split('.')[0]}')

        # SAFELY CLOSE APPLICATION IF BUILD FAILED
        if ret != 0:
            exit(ret)

        # BUILD FUNCTION DIRECTORY
        ret: int = mkdir(f'{self.active_directory}\\build\\build_{os.path.basename(self.path).split('.')[0]}'
                         f'\\functions')

        # SAFELY CLOSE APPLICATION IF BUILD FAILED
        if ret != 0:
            exit(ret)

        self.inheritance_map: dict = {}

    def get_type(self, var: str, line_num):
        var = var.strip()
        if len(var) != 0:
            match var[0]:
                case '\'':  # LINE STARTS WITH CHAR TYPE QUOTE
                    if len(var.replace('\'', '')) > 1:
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
                            f'Value Warning:{RESET} {ITALIC}\'{var}\'{RESET} '
                            f'is not a single character, {FG_BRIGHT_BLUE}line {line_num}{RESET}')
                    return 'char'
                case '\"':  # LINE STARTS WITH STRING TYPE QUOTE
                    return 'str'
                case '[':  # LINE STARTS WITH LIST DECLARATION
                    return f'list > {self.get_type(var[1:].split(',')[0], line_num)}'
                case '{':  # LINE STARTS WITH DICTIONARY DECLARATION
                    return (f'dict > {self.get_type(var[1:].split(',')[0].split(':')[0], line_num)} > '
                            f'{self.get_type(var[1:].split(',')[0].split(':')[1], line_num)}')
                case default:
                    if var.startswith('if') or var == 'true' or var == 'false':
                        # LINE IS BOOLEAN OR BOOLEAN EXPRESSION
                        return 'bool'
                    if var.isdigit():  # LINE IS INTEGER
                        return 'int'
                    try:
                        float(var)
                        return 'float'
                    except ValueError:
                        pass
        return 'void'

    def run(self):
        if self.ext == '.stream':  # CORRECT FILE FOR COMPILATION IF FILE IS A HIGH LEVEL STREAM APPLICATION
            self.file = self.correct_safety()

            print(
                f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {FG_BRIGHT_GREEN}'
                f'Successfully corrected file safety{RESET}')

            # WRITE STREAM LITE FILE
            write_file(self.output + '.sl', self.file)

        # GENERATE FUNCTIONAL STREAM LITE FILE
        self.file = self.generate_functional_file()

        print(
            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {FG_BRIGHT_GREEN}'
            f'Successfully generated functional stream lite{RESET}')

        # WRITE FUNCTIONAL STREAM LITE FILE
        write_file(self.output + '.fsl', self.file)

        print(
            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {FG_BRIGHT_GREEN}Successfully finished compilation'
            f'{RESET}')

    def correct_safety(self):
        new_file: str = ''

        namespace = None

        i = 0
        for line in self.file.split('\n'):
            i += 1
            # SPLIT MULTI ACTION LINES INTO SUB LINES
            split_line = line.split(';')

            for sub_line in split_line:
                # REMOVE COMMENTS AND REDUNDANT WHITESPACES
                new_line = sub_line.split('#')[0].strip()

                if new_line.startswith('if') or new_line.startswith('elif') or new_line.startswith('while') \
                        or new_line.startswith('else'):
                    if not new_line.__contains__('{'):
                        new_line += '{'
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}No end braces:'
                            f'{RESET} {ITALIC}\'{line}\'{RESET}, {FG_BRIGHT_BLUE}line {i}{RESET}')

                elif new_line.startswith('def'):
                    # ENSURE LINE CORRECTLY ENDS WITH AN ENDING TAG
                    if not new_line.__contains__('{'):
                        new_line += '{'
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}No end braces:'
                            f'{RESET} {ITALIC}\'{line}\'{RESET}, {FG_BRIGHT_BLUE}line {i}{RESET}')
                    # ENSURE LINE CONTAINS A RETURN TYPE
                    if not new_line.__contains__('->'):
                        # INSERT VOID RETURN
                        new_line = f'{new_line.split('(')[0].strip()} -> void ({new_line.split('(')[1].strip()}'
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
                            f'No function return notation:{RESET} {ITALIC}\'{line}\'{RESET}'
                            f', {FG_BRIGHT_BLUE}line {i}{RESET}')
                    # UPDATE NAME WITH NAMESPACE TAG
                    if namespace is not None:
                        new_line = f'def {namespace}.{line[4:]}'

                elif new_line.startswith('return'):
                    if len(new_line.split(' ')) != 2:
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
                            f'Return without statement:{RESET} {ITALIC}\'{line}\'{RESET}'
                            f', {FG_BRIGHT_BLUE}line {i}{RESET}')
                        new_line = f'{new_line.strip()} void'
                    if ((new_line.strip().endswith(')') and new_line.__contains__('('))
                            and not new_line.__contains__('exec')):
                        new_line = f'return exec {new_line[7:]}'
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
                            f'No exec tag:{RESET} {ITALIC}\'{line}\'{RESET}'
                            f', {FG_BRIGHT_BLUE}line {i}{RESET}')
                # SET NAMESPACE
                elif new_line.startswith('namespace'):
                    namespace = new_line[10:]
                    if namespace == 'end':
                        namespace = None
                    new_line = ''
                # CLOSE NAMESPACE
                elif new_line.startswith('end'):
                    namespace = None
                    new_line = ''
                # IF LINE CONTAINS = THAT IS NOT A PART OF A ==
                if re.search(r"(?<!=)=(?!=)", new_line):
                    # IF LINE DOES NOT HAVE AN EXEC TAG
                    if ((new_line.strip().endswith(')') and new_line.__contains__('('))
                            and not new_line.__contains__('exec')):
                        new_line = f'{new_line.split('=')[0].strip()} = exec {new_line.split('=', 1)[1].strip()}'
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
                            f'No exec tag:{RESET} {ITALIC}\'{line}\'{RESET}'
                            f', {FG_BRIGHT_BLUE}line {i}{RESET}')
                    # IF LINE DOES NOT HAVE A LET DECLARATION INSERT LET
                    if not new_line.startswith('let'):
                        new_line = f'let {new_line}'
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
                            f'No let tag:{RESET} {ITALIC}\'{line}\'{RESET}, {FG_BRIGHT_BLUE}line {i}{RESET}')
                    # IF LINE DOES NOT HAVE A TYPE ANNOTATION GENERATE TYPE NOTATION
                    if not new_line.split('=', 1)[0].__contains__(':'):
                        split_eq = new_line.split('=', 1)
                        var = split_eq[1].strip()
                        hand = split_eq[0]

                        type_declaration = self.get_type(var, i)

                        new_line = f'{hand.rstrip()}: {type_declaration} = {var}'

                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
                            f'No variable type notation:{RESET} {ITALIC}\'{line}\'{RESET}'
                            f', {FG_BRIGHT_BLUE}line {i}{RESET}')
                if new_line.startswith('let'):
                    # IF LINE HAS A DECLARATION BUT NOT A = STATEMENT INSERT VOID
                    if not re.search(r"(?<!=)=(?!=)", new_line):
                        new_line = f'{new_line.rstrip()} = void'
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
                            f'No equals statement:{RESET} {ITALIC}\'{line}\'{RESET}'
                            f', {FG_BRIGHT_BLUE}line {i}{RESET}')
                    elif new_line.strip().endswith('='):
                        new_line = f'{new_line.strip()} void'
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
                            f'No equals statement:{RESET} {ITALIC}\'{line}\'{RESET}'
                            f', {FG_BRIGHT_BLUE}line {i}{RESET}')
                    if line[-1] in ['*', '/', '+', '-']:
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
                            f'Unclosed operation:{RESET} {ITALIC}\'{line}\'{RESET}'
                            f', {FG_BRIGHT_BLUE}line {i}{RESET}')
                    # INSERT NAMESPACE TAG IN VARIABLE NAME
                    if namespace is not None:
                        new_line = f'let {namespace}.{new_line[4:]}'
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
                            f'No let tag:{RESET} {ITALIC}\'{line}\'{RESET}'
                            f', {FG_BRIGHT_BLUE}line {i}{RESET}')
                # IF LINE EXECUTES FUNCTION BUT DOES NOT CONTAIN AN EQUALS STATEMENT INSERT VOID STATEMENT
                if new_line.__contains__('exec') and not re.search(r"(?<!=)=(?!=)", new_line) and not line.startswith('return'):
                    new_line = f'let _:void = {new_line}'
                    print(
                        f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
                        f'No let statement:{RESET} {ITALIC}\'{line}\'{RESET}'
                        f', {FG_BRIGHT_BLUE}line {i}{RESET}')
                # INSERT LINE
                if new_line != '':
                    new_file += f'{new_line}\n'

        return new_file

    def generate_functional_file(self):
        new_file: str = ''

        in_foo = False

        split_file = self.file.split('\n')

        index = 0

        in_while = False
        while_var = ''

        while_jump_num = []

        general_depth = 0

        # LOOP THROUGH LINES IN FILE
        while index < len(split_file):
            # RETRIEVE LINE INFO
            new_line, line = split_file[index], split_file[index]

            # IF LINE IS AN IMPORT
            if line.startswith('from'):
                split_line = line[5:].split('import')
                directory = f'{self.active_directory}/{split_line[0].strip()[1:-1]}'
                new_line = ''
                # CHECK IF DIRECTORY EXISTS
                if os.path.exists(directory):
                    functions = split_line[1].split(',')
                    # READ PACKAGE DATA
                    with open(f'{directory}/function_manager.spk', 'r') as f:
                        pk_file = f.read()
                    pk_file = pk_file.split('\n')

                    # SPLIT PACKAGE DATA BY FUNCTION
                    split_pk_file = {line.split('<%$>')[0].strip(): line.split('<%$>')[1].strip() for line in pk_file}

                    # FOR FUNCTION IN IMPORT DECLARATION
                    for function in functions:
                        function = function.strip()
                        # IMPORT ALL FUNCTION IN THE PACKAGE MANAGER
                        if function == '*':
                            # LOOP THROUGH ALL FUNCTIONS IN PACKAGE MANAGER
                            for value in split_pk_file:
                                # BUILD PROTOTYPE FOR LINE
                                new_line += (f'prototype {value.split('.')[0]} \"{directory}/{value}\" -> '
                                             f'{split_pk_file[value]}\n')
                        else:
                            # PROTOTYPE FUNCTION
                            new_line += (f'prototype {function.split('.')[0]} \"{directory}/{function}\" -> '
                                         f'{split_pk_file[function]}\n')
                else:
                    print(
                        f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: '
                        f'{BOLD}{FG_BRIGHT_RED}No such directory:{RESET} '
                        f'Unable to import functions at \'{self.active_directory}/{directory}\'')
                    if self.syntactic:
                        exit(1)

                # REMOVE REDUNDANT LINES
                new_line = new_line.rstrip()

            elif line.startswith('def'):
                in_foo = True
                # GENERATE PROTOTYPE STATEMENT
                print(line.split('(')[1].split(')')[0].split(','), "SSS")
                new_line = (
                    f'prototype {line[4:].split(' ')[0]} \"{self.active_directory}\\build\\build_'
                    f'{os.path.basename(self.path).split('.')[0]}\\functions\\{line[4:].split(' ')[0]}.fsl\"'
                    f' -> {line.split('->')[1].strip().split(' ')[0]} '
                    f'{str(line.split('(')[1].split(')')[0].split(','))
                    .replace('\'', '').replace(']', ')')
                    .replace('[', '(')
                    if line.split('(')[1].split(')')[0].split(',') != [''] else '()'}')
                depth = 0
                data = ''
                jump_num = 1
                # LOAD FUNCTION CONTENTS AND JUMP OVER FUNCTION DATA TO PREVENT UNINTENDED EXECUTION
                for sub_line in split_file[index + 1:]:
                    if sub_line.startswith('if') or sub_line.startswith('else') or sub_line.startswith(
                            'elif') or sub_line.startswith('while'):
                        depth += 1
                    elif sub_line == '}':
                        depth -= 1
                        if depth <= 0:
                            break
                    jump_num += 1
                    data += sub_line + '\n'
                index += jump_num
                # WRITE FUNCTION DATA TO STREAM LITE FILE
                write_file(
                    f'{self.active_directory}\\build\\build_{os.path.basename(self.path).split('.')[0]}\\functions\\'
                    f'{line[4:].split(' ')[0]}.sl',
                    data)
                # COMPILE STREAM LITE
                sub_app = App(
                    f'{self.active_directory}\\build\\build_{os.path.basename(self.path).split('.')[0]}\\functions\\'
                    f'{line[4:].split(' ')[0]}.sl',
                    output=f'{self.active_directory}\\build\\build_{os.path.basename(self.path).split('.')[0]}'
                           f'\\functions\\{line[4:].split(' ')[0]}')
                sub_app.start()

            elif line.startswith('if'):
                general_depth += 1
                # BUILD CHECK
                new_file += f'let check: bool = {line[:-1]}'
                depth = 0
                jump_num = 0
                # CALCULATE DISTANCE TO STATEMENT END
                for sub_line in split_file[index:]:
                    if sub_line.startswith('if') or sub_line.startswith('else') or sub_line.startswith(
                            'elif') or sub_line.startswith('while'):
                        depth += 1
                    if sub_line == '}':
                        depth -= 1
                        if depth <= 0:
                            break
                    jump_num += 1
                # INSERT JUMP CHECK TO STATEMENT END (IF CONDITIONAL IS NOT TRUE IGNORE CODE)
                new_line = f'\ndo check {jump_num - 1}'

            elif line.startswith('else'):
                general_depth += 1
                # REVERSE CHECK (CHECK SHOULD HAVE BEEN BUILD IN PREVIOUS STATEMENTS)
                new_file += f'let check: bool = if not check:bool == true:bool'
                depth = 0
                jump_num = 0
                # CALCULATE DISTANCE TO STATEMENT END
                for sub_line in split_file[index:]:
                    if sub_line.startswith('if') or sub_line.startswith('else') or sub_line.startswith(
                            'elif') or sub_line.startswith('while'):
                        depth += 1
                    if sub_line == '}':
                        depth -= 1
                        if depth <= 0:
                            break
                    jump_num += 1
                # INSERT JUMP CHECK TO STATEMENT END (IF CONDITIONAL IS NOT TRUE IGNORE CODE)
                new_line = f'\ndo check {jump_num - 1}'

            elif line.startswith('elif'):
                general_depth += 1
                # FLIP CHECK (SHOULD HAVE BEEN GENERATED IN PREVIOUS STATEMENT) AND CHECK CONDITIONAL
                new_file += f'let check: bool = if not check:bool == true:bool and {line[5:-1]}'
                jump_num = 0
                depth = 0
                # CALCULATE DISTANCE TO STATEMENT END
                for sub_line in split_file[index:]:
                    if sub_line.startswith('if') or sub_line.startswith('else') or sub_line.startswith(
                            'elif') or sub_line.startswith('while'):
                        depth += 1
                    if sub_line == '}':
                        depth -= 1
                        if depth <= 0:
                            break
                    jump_num += 1
                # INSERT JUMP CHECK TO STATEMENT END (IF CONDITIONAL IS NOT TRUE IGNORE CODE)
                new_line = f'\ndo check {jump_num - 1}'

            elif line.startswith('while'):
                split_line = line[6:].split(' ')
                depth = 0
                jump_num = 0
                # CALCULATE DISTANCE TO STATEMENT END
                for sub_line in split_file[index:]:
                    if sub_line.startswith('if') or sub_line.startswith('else') or sub_line.startswith(
                            'elif') or sub_line.startswith('while'):
                        depth += 1
                    if sub_line == '}':
                        depth -= 1
                        if depth <= 0:
                            break
                    jump_num += 1
                # JUMP OVER LINE IF CONDITIONAL IS FALSE
                new_line = f'do {split_line[0]} {jump_num}'
                # INSERT JUMP NUM ONTO STACK (FOR NESTED WHILE LOOPS)
                while_jump_num.append(jump_num)
                in_while = True
                while_var = split_line[0]
                # RESET GENERAL DEPTH
                # general_depth = 0

            elif line == '}':
                general_depth -= 1
                if in_foo:
                    # INSERT RETURN STATEMENT IF NO RETURN STATEMENT EXISTS
                    if index - 1 >= 0 and not split_file[index - 1].startswith('return'):
                        new_file += f'return void\n'
                elif in_while and general_depth == 0:
                    # INSERT JUMP TO WHILE LOOP START IF CONDITIONAL STILL TRUE
                    new_line = f'do {while_var} {-while_jump_num.pop()}'
                else:
                    new_line = ''
                in_foo = False

            if new_line != '':
                new_file += f'{new_line}\n'

            index += 1

        return new_file


def main():
    # ENSURE ENOUGH ARGUMENTS ARE PASSED
    if len(sys.argv) < 1:
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Argument error:{RESET} '
              f'Not enough arguments {FG_BRIGHT_BLUE}(Must have 3: file path, return type, parameters){RESET}')
        sys.exit(1)
    elif len(sys.argv) > 1:
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
              f'Too many arguments{RESET} {FG_BRIGHT_BLUE}(Should have 3: file path, return type, parameters){RESET}')

    path = sys.argv[1]

    print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {FG_BRIGHT_GREEN}Running with parameters {RESET}:'
          f' \n{FG_BRIGHT_BLUE}path={path}{RESET}\n')

    app = App(path)
    app.start()


if __name__ == '__main__':
    main()
