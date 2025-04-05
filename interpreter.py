import importlib.util
import json
import math
import os
import sys
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


def fsl(file: str, type: str, param: dict, verbose):
    app = App(file, param, verbose)
    return app.run(type)


def load_module_from_path(filepath: str):
    """Loads a Python module from a given file path."""
    try:
        module_name = os.path.splitext(os.path.basename(filepath))[0]
        spec = importlib.util.spec_from_file_location(module_name, filepath)
        module = importlib.util.module_from_spec(spec)
        sys.modules[module_name] = module
        spec.loader.exec_module(module)
        return module
    except FileNotFoundError:
        print(f"Error: File not found at '{filepath}'")
        return None
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
        return None


def py(filepath: str, type: str, param: dict, verbose):
    """Loads a module and calls its 'main' function."""
    module = load_module_from_path(filepath)
    if module:
        try:
            main_func = getattr(module, 'main')
            return main_func(param)
        except AttributeError:
            return f"Error: 'main' function not found in '{filepath}'"
        except Exception as e:
            return f"An unexpected error occurred during main function call: {e}"
    else:
        return "Module could not be loaded"


def main(verbose=False):
    # ENSURE ENOUGH ARGUMENTS ARE PASSED
    if len(sys.argv) < 4:
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Argument error:{RESET} '
              f'Not enough arguments {FG_BRIGHT_BLUE}(Must have 3: file path, return type, parameters){RESET}')
        sys.exit(1)
    elif len(sys.argv) > 4:
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
              f'Too many arguments{RESET} {FG_BRIGHT_BLUE}(Should have 3: file path, return type, parameters){RESET}')

    f = sys.argv[1].replace('+', ' ')
    return_type = sys.argv[2]
    params = str(sys.argv[3]).replace('+', ' ')

    print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {FG_BRIGHT_GREEN}Running with parameters {RESET}:'
          f' \n{FG_BRIGHT_BLUE}return type={return_type}, parameters={params}{RESET}\n')

    try:
        params = json.loads(params)
    except Exception as e:
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}JSON error:{RESET} '
              f'Failed to load JSON parameters: {FG_BRIGHT_BLUE}\'{e}\'{RESET}')
        sys.exit(1)

    # GET FILE EXTENSION AND RUN PROPER PROCEDURES FOR THAT FILE EXTENSION
    _, ext = os.path.splitext(f)
    # SEARCH EXECUTION PIPELINE FOR PROPER FUNCTION HANDLER FOR FILE TYPE
    print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {FG_BRIGHT_GREEN}Successfully executed file{RESET}'
          f' at \'{f}\': \n{FG_BRIGHT_BLUE}{getattr(__import__(__name__), ext[1:])(f, return_type, params, verbose)}{RESET}')


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


def handler(f: str, return_type: str, params: dict, verbose):
    # GET FILE EXTENSION AND RUN PROPER PROCEDURES FOR THAT FILE EXTENSION
    _, ext = os.path.splitext(f)
    # SEARCH EXECUTION PIPELINE FOR PROPER FUNCTION HANDLER FOR FILE TYPE
    return getattr(__import__(__name__), ext[1:])(f, return_type, params, verbose)


def custom_split(expression):
    operators = ['**', '//', '+', '-', '*', '/', '~', '$', '@', '&']  # List of operators
    result = ['=']
    current = ""
    i = 0

    while i < len(expression):
        # CHECK FOR TWO CHAR OPERATORS (i.e. **, //)
        if expression[i:i + 2] in operators:
            if current.strip():  # INSERT NON OPERATOR
                result.append(current.strip())
            result.append(expression[i:i + 2])  # INSERT OPERATOR
            current = ""
            i += 2  # SKIP NEXT CHARACTER
        elif expression[i] in operators:  # HANDLE SINGLE LINE OPERATORS
            if current.strip():
                result.append(current.strip())  # INSERT NON OPERATOR
            result.append(expression[i])  # INSERT OPERATOR
            current = ""
            i += 1
        else:
            current += expression[i]  # INSERT NON OPERATOR CHARACTERS
            i += 1

    # ADD ANY REMAINING SECTIONS TO RESULT
    if current.strip():
        result.append(current.strip())

    return result


class App:
    def __init__(self, path: str, params: dict = None, verbose=False):
        # IF NO PARAMS VALUE SPECIFIED SET TO DEFAULT DICT
        self.verbose = verbose

        if params is None:
            params = {}

        self.path: str = os.path.abspath(path)
        self.params: dict = params

        # ENSURE FILE EXISTS
        if not os.path.exists(self.path):
            print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}File Not Found:{RESET} '
                  f'{ITALIC}\'{self.path}\'{RESET}')
            sys.exit(1)

        # READ FILE CONTENTS
        self.file: str = read_file(self.path)

        # SAFELY END PROGRAM IF READ FILE FAILED
        if self.file is None:
            sys.exit(1)

        self.calls = {}

    def check_conditional(self, conditional, local, line_num):
        parts = conditional[3:].split(' ')
        if parts == ['']:
            print(
                f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Conditional without body:'
                f' {ITALIC}\'{conditional}\'{RESET}, {FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
            sys.exit(1)
        final = False
        check_type = None
        check_obj = None
        upper_check = None
        upper_checks = ['and', 'or']
        tags = ['not']
        tag = None
        current = False
        for part in parts:
            # PART CONTAINS = (IS A CONDITION)
            if part.__contains__('=') or part.__contains__('<') or part.__contains__('>'):
                check_type = part
            # PART IS A CONNECTING CHECK BETWEEN CONDITIONALS
            elif part in upper_checks:
                upper_check = part
            # PART IS A TAG (i.e., not)
            elif part in tags:
                tag = part
            # PART IS NOT A CHECK TYPE AND CHECK TYPE IS NOT SET (PART IS A CHECK OBJECT)
            elif check_type is None and not part.__contains__('=') and not part.__contains__('<') and not part.__contains__('>'):
                try:
                    check_obj = self.evaluate(part.split(':')[0], part.split(':')[1], local, line_num)
                except IndexError:
                    print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}'
                          f'Typeless conditional object: {ITALIC}\'{part}\'{RESET}, {FG_BRIGHT_BLUE}line {line_num + 1}'
                          f'{RESET}')
                    sys.exit(1)
            # PART IS A CHECK CONDITION
            else:
                match check_type:
                    case '==':
                        try:
                            current = check_obj == self.evaluate(part.split(':')[0], part.split(':')[1], local,
                                                                 line_num)
                        except IndexError:
                            print(
                                f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}'
                                f'Typeless conditional check: {ITALIC}\'{part}\'{RESET}')
                            sys.exit(1)
                        check_type = None
                        check_obj = None
                    case '<=':
                        try:
                            current = check_obj <= self.evaluate(part.split(':')[0], part.split(':')[1], local,
                                                                 line_num)
                        except IndexError:
                            print(
                                f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}'
                                f'Typeless conditional check: {ITALIC}\'{part}\'{RESET}')
                            sys.exit(1)
                        check_type = None
                        check_obj = None
                    case '>=':
                        try:
                            current = check_obj >= self.evaluate(part.split(':')[0], part.split(':')[1], local,
                                                                 line_num)
                        except IndexError:
                            print(
                                f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}'
                                f'Typeless conditional check: {ITALIC}\'{part}\'{RESET}, '
                                f'{FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                            sys.exit(1)
                        check_type = None
                        check_obj = None
                    case '<':
                        try:
                            current = check_obj < self.evaluate(part.split(':')[0], part.split(':')[1], local,
                                                                 line_num)
                        except IndexError:
                            print(
                                f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}'
                                f'Typeless conditional check: {ITALIC}\'{part}\'{RESET}')
                            sys.exit(1)
                        check_type = None
                        check_obj = None
                    case '>':
                        try:
                            current = check_obj > self.evaluate(part.split(':')[0], part.split(':')[1], local,
                                                                 line_num)
                        except IndexError:
                            print(
                                f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}'
                                f'Typeless conditional check: {ITALIC}\'{part}\'{RESET}, '
                                f'{FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                            sys.exit(1)
                        check_type = None
                        check_obj = None
                if tag == 'not':
                    current = not current
                tag = None
                match upper_check:
                    case None:
                        final = current
                    case 'and':
                        final = final and current
                    case 'or':
                        final = final or current
                    case 'xor':
                        final = final ^ current
        return final

    def evaluate(self, value, var_type, local, line_num):
        if value == 'None':
            return None
        if value.startswith('exec'):
            param_build = {}
            index = 0

            if not self.calls.__contains__(value[5:].split('(')[0].strip()):
                print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}'
                      f'Function Not Found:{RESET} '
                      f'{ITALIC}\'{value[5:].split('(')[0]}\'{RESET}, registered calls: {self.calls}, {FG_BRIGHT_BLUE}'
                      f'line {line_num + 1}{RESET}')
                sys.exit(1)

            if value.split('(', 1)[1].rsplit(')', 1)[0].split(',') != ['']:
                for param in value.split('(', 1)[1].rsplit(')', 1)[0].split(','):
                    try:
                        key, param_type = self.calls[value[5:].split('(')[0]][2][index]
                    except IndexError:
                        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}'
                              f'Mismatched arguments:{RESET} '
                              f'{ITALIC}\'{value[5:].split('(')[0].strip()}\' takes in {FG_BRIGHT_BLUE}{str(self.calls[value[5:].split('(')[0]][2]).replace('[', '(').replace(']', ')')}{RESET}, given {FG_BRIGHT_BLUE}{str(value.split('(')[1].split(')')[0].split(',')).replace('[', '(').replace(']', ')')}{RESET} {RESET}, {FG_BRIGHT_BLUE}'
                              f'line {line_num + 1}{RESET}')
                        sys.exit(1)
                    param_build[key.strip()] = self.evaluate(param.strip(), param_type.strip(), local, line_num)
                    index += 1

            result = handler(self.calls[value[5:].split('(')[0]][0], self.calls[value[5:].split('(')[0]][1],
                             param_build | local, self.verbose)

            if self.verbose:
                print(
                    f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {FG_BRIGHT_GREEN}'
                    f'Successfully executed function file{RESET}'
                    f' at \'{self.calls[value[5:].split('(')[0]][0]}\': {FG_BRIGHT_BLUE}{result}, '
                    f'{FG_BRIGHT_BLUE}line {line_num + 1}{RESET}\n')

            return result
        if value.__contains__('.') and not value.split('.')[0].isdigit():
            try:
                obj = local[value.split('.')[0]]
                try:
                    if type(obj) is dict:
                        value = obj[self.evaluate(value.split('.')[1].split(':')[0], value.split('.')[1].split(':')[1],
                                                  local, line_num)]
                    elif type(obj) is list:
                        value = obj[self.evaluate(value.split('.')[1].split(':')[0], value.split('.')[1].split(':')[1],
                                                  local, line_num)]
                    else:
                        value = getattr(obj, self.evaluate(value.split('.')[1].split(':')[0], value.split('.')[1].split(':')[1],
                                                  local, line_num))
                except IndexError:
                    try:
                        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}'
                              f'Index Error:{RESET} {ITALIC}{obj}.'
                              f'{self.evaluate(value.split('.')[1].split(':')[0], value.split('.')[1].split(':')[1], 
                                               local, line_num)} '
                              f'out of bounds{RESET}, {FG_BRIGHT_BLUE}(length = {len(obj)}){RESET}, {FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                    except:
                        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}'
                              f'Syntax Error:{RESET} Statement does not have type annotations: {ITALIC}{value}'
                              f'{RESET}, {FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                    sys.exit(1)
                except KeyError:
                    print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Key Error:{RESET} '
                          f'\'{value}\' does not exist {FG_BRIGHT_BLUE}(Target not found: \'{value.split('.')[1]}\' in '
                          f'{local[value.split('.')[0]]}){RESET}'
                          f', {FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                    sys.exit(1)
            except KeyError:
                print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Key Error:{RESET} '
                      f'\'{value}\' does not exist {FG_BRIGHT_BLUE}(Target not found: \'{value.split('.')[0]}\' in '
                      f'{local}){RESET}, {FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                sys.exit(1)
            return value
        if value in local:
            return local[value]
        match var_type:
            case 'void':
                return None
            case 'obj':
                return value
            case 'int':
                try:
                    return int(value)
                except ValueError:
                    print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Value Error:{RESET} '
                          f'{ITALIC}\'{value}\'{RESET} is not an integer, {FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                    sys.exit(1)
            case 'float':
                try:
                    return float(value)
                except ValueError:
                    print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Value Error:{RESET} '
                          f'{ITALIC}\'{value}\'{RESET} is not a float, {FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                    sys.exit(1)
            case 'str':
                return value.replace('"', '')
            case 'char':
                value = value.replace('\'', '')
                return value
            case 'bool':
                if value.startswith('if'):
                    return self.check_conditional(value, local, line_num)
                else:
                    try:
                        return True if value.capitalize() == 'True' else False
                    except ValueError:
                        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Value Error:'
                              f'{RESET} {ITALIC}\'{value}\'{RESET} is not a boolean, '
                              f'{FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                        sys.exit(1)
        if var_type.__contains__('>'):
            true_type = var_type.split('>')[0].strip()
            if true_type == 'list':
                return [self.evaluate(item, var_type.split('>')[1].strip(), local, line_num) for item in value
                        .replace('[', '').replace(']', '').split(',')]
            elif true_type == 'dict':
                return {
                    self.evaluate(item.split(':')[0].strip(), var_type.split('>')[1].strip(), local, line_num):
                        self.evaluate(item.split(':')[1].strip(), var_type.split('>')[2].strip(), local, line_num)
                    for item in value[1: -1].split(',')}
        return None

    def evaluate_multi(self, value, var_type, local, line_num):
        # SPLIT VALUE BY OPERATORS
        result = custom_split(value)

        match var_type:
            case 'int':
                final = 0
            case 'list':
                final = []
            case 'dict':
                final = {}
            case 'str':
                final = ''
            case default:
                final = None

        operations = ['+', '-', '*', '/', '//', '**', '=', '%', '@', '&']

        current_operation = ''

        for item in result:
            if item in operations:
                current_operation = item
                continue
            if item == '~':
                final = math.floor(final)
                continue
            if item == '$':
                try:
                    final = len(final)
                except TypeError:
                    print(
                        f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: '
                        f'{FG_BRIGHT_RED}Invalid opperation:{RESET}'
                        f' cannot take the length of {FG_BRIGHT_BLUE}{final}{RESET}, invalid type ({type(final)})'
                        f', {FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                    sys.exit(1)
                continue
            if var_type == 'int':
                true_item = self.evaluate(item, 'int' if not item.__contains__('.') else 'float', local, line_num)
                match current_operation:
                    case '+':
                        final += true_item
                    case '-':
                        final -= true_item
                    case '*':
                        final *= true_item
                    case '/':
                        final /= true_item
                    case '//':
                        final **= (1 / true_item)
                    case '**':
                        final = math.pow(final, true_item)
                    case '=':
                        final = true_item
                    case '%':
                        final %= true_item
                    case default:
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: '
                            f'{FG_BRIGHT_RED}Invalid integer operation:{RESET}'
                            f':\n{FG_BRIGHT_BLUE}\'{current_operation}\', {value}{RESET}'
                            f', {FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                        sys.exit(1)
            elif var_type == 'str':
                match current_operation:
                    case '+':
                        final += self.evaluate(item, 'str', local, line_num)
                    case '*':
                        final *= self.evaluate(item, 'int', local, line_num)
                    case '=':
                        final = self.evaluate(item, 'str', local, line_num)
                    case default:
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: '
                            f'{FG_BRIGHT_RED}Invalid string operation:{RESET}'
                            f':\n{FG_BRIGHT_BLUE}\'{current_operation}\', {value}{RESET}'
                            f', {FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                        sys.exit(1)
            elif var_type == 'list':
                match current_operation:
                    case '+':
                        final = final + self.evaluate(item.split(':')[0], 'list', local, line_num)
                    case '=':
                        final = self.evaluate(item, 'list', local, line_num)
                    case '@':
                        final = final[self.evaluate(item, 'list', local, line_num):]
                    case '&':
                        final = final[:self.evaluate(item, 'list', local, line_num)]
                    case default:
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: '
                            f'{FG_BRIGHT_RED}Invalid string operation:{RESET}'
                            f':\n{FG_BRIGHT_BLUE}\'{current_operation}\', {value}{RESET}'
                            f', {FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                        sys.exit(1)
            elif var_type == 'dict':
                match current_operation:
                    case '+':
                        try:
                            final = final | self.evaluate(item.rsplit(':', 1)[0], item.rsplit(':', 1)[1], local,
                                                          line_num)
                        except TypeError:
                            print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: '
                                  f'{FG_BRIGHT_RED}TypeError:{RESET} '
                                  f'Invalid statement return \'{item}\' {FG_BRIGHT_BLUE}(Returns '
                                  f'\'{self.evaluate(item.rsplit(':', 1)[0], item.rsplit(':', 1)[1], 
                                                     local, line_num)}\''
                                  f' which cannot be merged with dictionary){RESET}'
                                  f', {FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                            exit(1)
                    case '=':
                        final = self.evaluate(item, 'dict', local, line_num)
                    case default:
                        print(
                            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: '
                            f'{FG_BRIGHT_RED}Invalid string operation:{RESET}'
                            f':\n{FG_BRIGHT_BLUE}\'{current_operation}\', {value}{RESET}'
                            f', {FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
                        sys.exit(1)

        return final

    def run(self, return_type):
        local = self.params
        # jump_num = 0
        split_file = self.file.split('\n')
        i = 0
        while i < len(split_file):
            line = split_file[i]
            if line.startswith('let'):
                sub_line = line.split('=')[1]
                if sub_line.__contains__('+') or sub_line.__contains__('-') or sub_line.__contains__('*') or \
                        sub_line.__contains__('/') or sub_line.__contains__('@') or sub_line.__contains__('&') or \
                        sub_line.__contains__('$') or sub_line.__contains__('~'):
                    key = line.split('=', 1)[0].split(':', 1)[0][4:].strip()
                    value = (self.evaluate_multi(line.split('=', 1)[1].strip(),
                                                      line.split('=', 1)[0].split(':', 1)[1].strip(), local, i))
                    if not key == '_':
                        local[key] = value
                else:
                    key = line.split('=', 1)[0].split(':', 1)[0][4:].strip()
                    value = (
                        self.evaluate(line.split('=', 1)[1].strip(),
                                      line.split('=', 1)[0].split(':', 1)[1].strip(), local, i))
                    if not key == '_':
                        local[key] = value
            elif line.startswith('!'):
                target = line[1:].strip()
                if target in local:
                    del local[target]
            elif line.startswith('do'):
                split_line = line[3:].split(' ')
                if not self.evaluate(split_line[0], 'bool', local, i):
                    i += self.evaluate(split_line[1], 'int', local, i)
            elif line.startswith('prototype'):
                split_line = line[10:].split(' ', 1)
                self.calls[split_line[0]] = [split_line[1].split('->')[0].strip()[1:-1].replace('+', ' '),
                                             split_line[1].split('->')[1].split('(')[0].strip(),
                                             [(param, param_type) for value in split_line[1].split('(')[1].split(')')[0]
                                             .split(',')
                                              for param, param_type in [value.split(':')]]
                                             if split_line[1].split('(')[1][:-1].split(',') != [''] else []]
            elif line.startswith('return'):
                if self.verbose:
                    print(
                        f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {FG_BRIGHT_GREEN}'
                        f'Program finished with local data{RESET}:'
                        f'{FG_BRIGHT_BLUE} local data={local}{RESET}')
                return self.evaluate(line[7:].strip(), return_type, local, i)
            i += 1
        print(
            f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
            f'No return statement{RESET}')
        if self.verbose:
            print(
                f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {FG_BRIGHT_GREEN}'
                f'Program finished with local data{RESET}:'
                f'{FG_BRIGHT_BLUE} local data={local}{RESET}')
        return None


if __name__ == '__main__':
    # EXAMPLE EXECUTION: python interpreter.py build\build_test2\test_sub_func.fsl int '{\"asdasdadsda\": false}'
    main()
