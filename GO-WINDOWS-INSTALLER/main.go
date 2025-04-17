package main

import (
	"fmt"
	_ "fmt"
	"io/ioutil"
	_ "io/ioutil"
	"os"
	_ "os"
	"os/exec"
	_ "os/exec"
	"path/filepath"
	_ "path/filepath"
	"runtime"
	_ "runtime"
	"strings"
	_ "strings"
)

const pythonCode = `
import importlib.util
import json
import math
import os
import sys
import time


# Define print variables
RESET = "\033[0m"
BOLD = "\033[1m"
ITALIC = "\033[3m"

FG_BRIGHT_RED = "\033[91m"
FG_BRIGHT_GREEN = "\033[92m"
FG_BRIGHT_YELLOW = "\033[93m"
FG_BRIGHT_BLUE = "\033[94m"
FG_BRIGHT_CYAN = "\033[96m"

start_time = time.time()


def fsl(file: str, return_type: str, param: dict, verbose):
    app = App(file, param, verbose)
    return app.run(return_type)


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


def py(filepath: str, _return_type: str, param: dict, _verbose):
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
    if len(sys.argv) < 5:
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Argument error:{RESET} '
              f'Not enough arguments {FG_BRIGHT_BLUE}(Must have 3: file path, return type, parameters){RESET}')
        sys.exit(1)
    elif len(sys.argv) > 5:
        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_YELLOW}'
              f'Too many arguments{RESET} {FG_BRIGHT_BLUE}(Should have 3: file path, return type, parameters){RESET}')

    f = sys.argv[1]
    return_type = sys.argv[2]
    params = str(sys.argv[3])
    verbose = str(sys.argv[4]).lower()
    if verbose == "true":
        verbose = True
    else:
        verbose = False

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
          f' at \'{f}\': \n{FG_BRIGHT_BLUE}{getattr(__import__(__name__), ext[1:])(f, return_type, params, verbose)}'
          f'{RESET}')


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
    operators = ['**', '//', '+', '-', '*', '/', '~', '$', '@', '&', '%']  # List of operators
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
            elif check_type is None and not part.__contains__('=') and not part.__contains__('<') \
                    and not part.__contains__('>'):
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
                              f'{ITALIC}\'{value[5:].split('(')[0].strip()}\' takes in {FG_BRIGHT_BLUE}'
                              f'{str(self.calls[value[5:].split('(')[0]][2]).replace('[', '(').replace(']', ')')}'
                              f'{RESET}, given {FG_BRIGHT_BLUE}{str(value.split('(')[1].split(')')[0].split(','))
                                                                .replace('[', '(').replace(']', ')')}'
                              f'{RESET} {RESET}, {FG_BRIGHT_BLUE}'
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
                    elif type(obj) is list or type(obj) is str:
                        value = obj[self.evaluate(value.split('.')[1].split(':')[0], value.split('.')[1].split(':')[1],
                                                  local, line_num)]
                    else:
                        value = getattr(obj, self.evaluate(value.split('.')[1].split(':')[0],
                                                           value.split('.')[1].split(':')[1], local, line_num))
                except IndexError:
                    try:
                        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}'
                              f'Index Error:{RESET} {ITALIC}{obj}.'
                              f'{self.evaluate(value.split('.')[1].split(':')[0], value.split('.')[1].split(':')[1], 
                                               local, line_num)} '
                              f'out of bounds{RESET}, {FG_BRIGHT_BLUE}(length = {len(obj)}){RESET}, '
                              f'{FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
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
                    print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Value Error:'
                          f'{RESET} {ITALIC}\'{value}\'{RESET} is not an integer, {FG_BRIGHT_BLUE}line {line_num + 1}'
                          f'{RESET}')
                    sys.exit(1)
            case 'float':
                try:
                    return float(value)
                except ValueError:
                    print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Value Error:'
                          f'{RESET} {ITALIC}\'{value}\'{RESET} is not a float, {FG_BRIGHT_BLUE}line {line_num + 1}'
                          f'{RESET}')
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

        final = None
        match var_type:
            case 'int':
                final = 0
            case 'list':
                final = []
            case 'dict':
                final = {}
            case 'str':
                final = ''

        operations = ['+', '-', '*', '/', '//', '**', '=', '%', '@', '&', '%']

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
                        f'{FG_BRIGHT_RED}Invalid operation:{RESET}'
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
                            f':\n{FG_BRIGHT_BLUE}\'{default}\', {value}{RESET}'
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
                            f':\n{FG_BRIGHT_BLUE}\'{default}\', {value}{RESET}'
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
                            f':\n{FG_BRIGHT_BLUE}\'{default}\', {value}{RESET}'
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
                            f':\n{FG_BRIGHT_BLUE}\'{default}\', {value}{RESET}'
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
                        sub_line.__contains__('$') or sub_line.__contains__('~') or sub_line.__contains__('%'):
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
    main()
`

const toml_code = `
[package]
name = "Rust-Stream-Compiler"
version = "0.1.0"
edition = "2021"

[dependencies]
serde_json = "1.0.140"
libloading = "0.8.6"
lazy_static = "1.5.0"
regex = "1.11.1"
time = "0.3.41"
csv = "1.3.1"
pyo3 = "0.24.1"
windows = { version = "0.61.1", features = ["Win32_Foundation", "Win32_UI_WindowsAndMessaging", "Win32_Graphics_Gdi", "Win32_System", "Win32_System_LibraryLoader"] }
`

const rs_code = `
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
            format!("{}\\build\\build_{}\\{}",
                    active_directory, basename, basename)
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
            write_file(&format!("{}.sl", self.output.as_ref().unwrap()), &self.file).expect("Failed to write .sl file");
        }

        // GENERATE FUNCTIONAL STREAM LITE FILE
        self.file = self.generate_functional_file();

        println!("{}{:?}{}: {}Successfully generated functional stream lite{}",
                 FG_BRIGHT_CYAN, get_elapsed_time(), RESET,
                 FG_BRIGHT_GREEN, RESET);

        // WRITE FUNCTIONAL STREAM LITE FILE
        write_file(&format!("{}.fsl", self.output.as_ref().unwrap()), &self.file).expect("Failed to write .fsl file");

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
                // r"([^=]|^)=([^=]|$)"
                let re = Regex::new(r"([^=<>]|^)=([^=<>]|$)").unwrap();
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
                    !new_line.starts_with("prototype") && !new_line.starts_with("def") {
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
                if new_line.contains("exec") && !has_lone_equals(&new_line) && !new_line.starts_with("return") {
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
                                    if !part.contains('=') && !part.contains(">") && !part.contains("<") && !upper_checks.contains(&part) &&
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
                let function_name = line[4..].split(' ').next().unwrap_or("");
                let base_name = Path::new(&self.path)
                    .file_name()
                    .unwrap()
                    .to_str()
                    .unwrap()
                    .split('.')
                    .collect::<Vec<&str>>()[0];
                let dest_type = line.split("->").collect::<Vec<&str>>()[1].trim().split(' ').collect::<Vec<&str>>()[0];
                new_line = format!("prototype {} \"{}\\build\\build_{}\\functions\\{}.fsl\" -> {} ({})",
                                   function_name,  self.active_directory, base_name, function_name, dest_type, if let Some(open_paren_index) = line.find('(') { let close_paren_index = line.find(')').unwrap(); let parameter_str = &line[open_paren_index + 1..close_paren_index]; if parameter_str.is_empty() { "()".to_string() } else { parameter_str.split(',').collect::<Vec<&str>>().join(", ").replace('\'', "").replace('[', "(").replace(']', ")") } } else { "()".to_string() });
                let mut depth = 0;
                let mut data: String = format!("{}\n", new_line.clone());
                let mut jump_num = 1;
                for sub_line in split_file.iter().skip(index + 1) {
                    if sub_line.starts_with("if") || sub_line.starts_with("else") || sub_line.starts_with("elif") || sub_line.starts_with("while") {
                        depth += 1;
                    }
                    else if *sub_line == "}" {
                        depth -= 1;
                        if depth < 0 {
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
                            // let parts: Vec<&str> = self.path.split(std::path::MAIN_SEPARATOR).collect();
                            // let split_part: Vec<&str> = parts.last().map(|s| s.to_string()).unwrap_or_else(String::new).split(".").collect();
                            // split_part.first().map(|s| s.to_string()).unwrap_or_else(String::new)
                            let parts: Vec<&str> = self.path.split(std::path::MAIN_SEPARATOR).collect();
                            let file_name = parts.last().map(|s| s.to_string()).unwrap_or_else(String::new);
                            let split_part: Vec<&str> = file_name.split(".").collect();
                            split_part.first().map(|s| s.to_string()).unwrap_or_else(String::new)
                        },
                        line[4..].split(' ').next().unwrap_or("")
                    ),
                    Some(format!(
                        "{}\\build\\build_{}\\functions\\{}",
                        self.active_directory,
                        {
                            // let parts: Vec<&str> = self.path.split(std::path::MAIN_SEPARATOR).collect();
                            // let split_part: Vec<&str> = parts.last().map(|s| s.to_string()).unwrap_or_else(String::new).split(".").collect();
                            // split_part.first().map(|s| s.to_string()).unwrap_or_else(String::new)
                            let parts: Vec<&str> = self.path.split(std::path::MAIN_SEPARATOR).collect();
                            let file_name = parts.last().map(|s| s.to_string()).unwrap_or_else(String::new);
                            let split_part: Vec<&str> = file_name.split(".").collect();
                            split_part.first().map(|s| s.to_string()).unwrap_or_else(String::new)
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
                new_file += &format!("let check: bool = if {}\n", line[3..line.len() - 1].trim());
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
                new_file += "let check: bool = if not check:bool == true:bool\n";
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
                new_file += &format!("let check: bool = if not check:bool == true:bool and {}\n", line[5..line.len() - 1].trim());
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
                    if sub_line.starts_with("while") {
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
                    // let parts: Vec<&str> = val.split("=").collect();
                    // val = match (parts.first()) {
                    //     Some(val) => {
                    //         val
                    //     }
                    //     None => {
                    //         ""
                    //     }
                    // };

                    if val.starts_with("let") {
                        val = &val[4..];
                    }
                    new_file += &format!(
                        "let {}: bool = {}\n",
                        val.split('=').next().unwrap().split(':').next().unwrap().trim(),
                        val.split('=').nth(1).unwrap().trim()
                    );
                    val = val.split('=').next().unwrap().split(':').next().unwrap().trim();
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
                if in_while && *while_general_depth.last().unwrap() == 0 {
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

    let args: Vec<String> = env::args().collect();

    if args.len() < 4 { // Rust args include the program name itself
        eprintln!(
            "{}{:?}{}: {}{}Argument error:{} Not enough arguments {}(Must have 3: file path, output path, syntactic mode){}",
            FG_BRIGHT_CYAN, get_elapsed_time(), RESET, BOLD, FG_BRIGHT_RED, RESET, FG_BRIGHT_BLUE, RESET
        );
        exit(1);
    } else if args.len() > 4 {
        eprintln!(
            "{}{:?}{}: {}Too many arguments{}{} {}(Should have 3: file path, output path, syntactic mode){}",
            FG_BRIGHT_CYAN, get_elapsed_time(), RESET, BOLD, FG_BRIGHT_YELLOW, RESET, FG_BRIGHT_BLUE, RESET
        );
    }

    let path: &str = &args[1].clone(); // Clone the path to own it.
    let output: Option<&str> = if args[2] == "none" {
        None
    } else {
        None
    };
    let syntactic: bool = if args.len() > 1 && args[1].to_lowercase() == "true" {
        true
    } else {
        false
    };

    println!(
        "{}{:?}{}: {}Running with parameters {}:{}{} \n{}path: {}, output path: {:?}, syntactic mode: {}{}",
        FG_BRIGHT_CYAN, get_elapsed_time(), RESET, FG_BRIGHT_GREEN, RESET, FG_BRIGHT_BLUE, RESET, FG_BRIGHT_BLUE, path,
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
`

const devkit_code = `
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
                os.system(f"./C:\\Stream\\bin\\stream-c.exe "
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
                os.system(f"./C:\\Stream\\bin\\stream-c.exe "
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
                   {blue}` + "`" + `-.|{reset}               {orange}"{reset}
'    {blue}` + "`" + `+_ +{reset}          {yellow}|{reset}
.'           {indigo}` + "`" + `x{reset}
.'                {indigo}={reset}
{green}O{reset}    {magenta}_.-'     ..{reset}
{green}( >-='{reset}            {magenta}` + "``" + `.._                  _.{reset}
{magenta}-'{reset}{green}  / \\"{reset}                   {magenta}` + "`" + `-a:f,-'-._    ,-'{reset}
{green}()){reset}
{magenta}_{reset} {green}b b{reset}       ____  __..__   __
{magenta}` + "`" + `"{reset}        |    |/ _||__|_/  |_
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
`

const web_code = `
import http.server
import importlib.util
import json
import math
import os
import sys
import time
import socketserver

PORT = 8000

# Define print variables
RESET = "\033[0m"
BOLD = "\033[1m"
ITALIC = "\033[3m"

FG_BRIGHT_RED = "\033[91m"
FG_BRIGHT_GREEN = "\033[92m"
FG_BRIGHT_YELLOW = "\033[93m"
FG_BRIGHT_BLUE = "\033[94m"
FG_BRIGHT_CYAN = "\033[96m"

start_time = time.time()


def fsl(file: str, return_type: str, param: dict, verbose):
    app = App(file, param, verbose)
    return app.run(return_type)


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


def py(filepath: str, _return_type: str, param: dict, _verbose):
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
    operators = ['**', '//', '+', '-', '*', '/', '~', '$', '@', '&', '%']  # List of operators
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
    def __init__(self, path: str = None, content: str = None, params: dict = None, verbose=False):
        # IF NO PARAMS VALUE SPECIFIED SET TO DEFAULT DICT
        self.verbose = verbose

        if params is None:
            params = {}

        self.params: dict = params
        if path is not None:
            self.path: str = os.path.abspath(path)

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
        elif content is not None:
            self.file = content
        else:
            print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Invalid parameters:{RESET} '
                  f'{ITALIC}no path or content provided{RESET}')
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
            elif check_type is None and not part.__contains__('=') and not part.__contains__('<') \
                    and not part.__contains__('>'):
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
                              f'{ITALIC}\'{value[5:].split('(')[0].strip()}\' takes in {FG_BRIGHT_BLUE}'
                              f'{str(self.calls[value[5:].split('(')[0]][2]).replace('[', '(').replace(']', ')')}'
                              f'{RESET}, given {FG_BRIGHT_BLUE}{str(value.split('(')[1].split(')')[0].split(','))
                                                                .replace('[', '(').replace(']', ')')}'
                              f'{RESET} {RESET}, {FG_BRIGHT_BLUE}'
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
                    elif type(obj) is list or type(obj) is str:
                        value = obj[self.evaluate(value.split('.')[1].split(':')[0], value.split('.')[1].split(':')[1],
                                                  local, line_num)]
                    else:
                        value = getattr(obj, self.evaluate(value.split('.')[1].split(':')[0],
                                                           value.split('.')[1].split(':')[1], local, line_num))
                except IndexError:
                    try:
                        print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}'
                              f'Index Error:{RESET} {ITALIC}{obj}.'
                              f'{self.evaluate(value.split('.')[1].split(':')[0], value.split('.')[1].split(':')[1], 
                                               local, line_num)} '
                              f'out of bounds{RESET}, {FG_BRIGHT_BLUE}(length = {len(obj)}){RESET}, '
                              f'{FG_BRIGHT_BLUE}line {line_num + 1}{RESET}')
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
                    print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Value Error:'
                          f'{RESET} {ITALIC}\'{value}\'{RESET} is not an integer, {FG_BRIGHT_BLUE}line {line_num + 1}'
                          f'{RESET}')
                    sys.exit(1)
            case 'float':
                try:
                    return float(value)
                except ValueError:
                    print(f'{FG_BRIGHT_CYAN}{time.time() - start_time}{RESET}: {BOLD}{FG_BRIGHT_RED}Value Error:'
                          f'{RESET} {ITALIC}\'{value}\'{RESET} is not a float, {FG_BRIGHT_BLUE}line {line_num + 1}'
                          f'{RESET}')
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

        final = None
        match var_type:
            case 'int':
                final = 0
            case 'list':
                final = []
            case 'dict':
                final = {}
            case 'str':
                final = ''

        operations = ['+', '-', '*', '/', '//', '**', '=', '%', '@', '&', '%']

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
                        f'{FG_BRIGHT_RED}Invalid operation:{RESET}'
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
                            f':\n{FG_BRIGHT_BLUE}\'{default}\', {value}{RESET}'
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
                            f':\n{FG_BRIGHT_BLUE}\'{default}\', {value}{RESET}'
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
                            f':\n{FG_BRIGHT_BLUE}\'{default}\', {value}{RESET}'
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
                            f':\n{FG_BRIGHT_BLUE}\'{default}\', {value}{RESET}'
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
                        sub_line.__contains__('$') or sub_line.__contains__('~') or sub_line.__contains__('%'):
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


class MyHandler(http.server.SimpleHTTPRequestHandler):
    def push_data(self, content: str):
        new_content = ""

        split_content = content.split("{")

        for section in split_content:
            if not section.__contains__("}"):
                new_content += section
            else:
                value = section.split("}")[0]
                etc = section.split("}")[1]

                ret = value.split("]")[0].split("[")[1]

                value = value.split("]")[1]

                new_value = App(content=value).run(ret)

                new_content += str(new_value) + etc

        return new_content

    def do_GET(self):
        self.path = self.path.replace("/", "")
        if self.path == '':
            self.path = "index"
        try:
            # Read file as string instead of bytes
            with open(f'{self.path}.html', 'r', encoding='utf-8') as file:
                content = file.read()

            # Process the content
            processed_content = self.push_data(content)

            # Convert to bytes only at the final step
            content_bytes = processed_content.encode('utf-8')

            self.send_response(200)
            self.send_header('Content-type', 'text/html; charset=utf-8')
            self.send_header('Content-length', len(content_bytes))
            self.end_headers()
            self.wfile.write(content_bytes)
            return
        except FileNotFoundError:
            self.send_error(404, "File not found")
            return
        return http.server.SimpleHTTPRequestHandler.do_GET(self)


with socketserver.TCPServer(("", PORT), MyHandler) as httpd:
    print(f"server live at 127.0.0.1:{PORT}")
    httpd.serve_forever()
`

const web_compiler_code = `
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
`

func cmd(command string, arguments []string, dir string) {
	cmd := exec.Command(command, arguments...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if dir != "" {
		cmd.Dir = dir
	}
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command '%s': %s\n", command, err)
		return
	}
}

func tempPythonFile(content string) string {
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "temp_script.py")
	err := ioutil.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error writing temporary Python file:", err)
		return ""
	}
	return tempFile
}

func mv(oldPath, newPath string) error {
	// Ensure the destination directory exists.
	destDir := filepath.Dir(newPath)
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating destination directory: %w", err) //wrap
	}

	// Use os.Rename to move the file.
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Errorf("error moving file: %w", err) //wrap
	}
	return nil
}

func mkdir(dirPath string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err) // Wrap for context
	}
	return nil
}

func addPath(dirPath string) error {
	// 1. Get the current PATH value.
	pathEnv := os.Getenv("PATH")

	// 2. Split the PATH into individual paths.
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	// 3. Check if the directory path is already in the PATH.
	for _, path := range paths {
		if path == dirPath {
			return nil // Already in PATH, nothing to do.
		}
	}

	// 4. If the path is not in the PATH, append it.
	newPath := pathEnv
	if runtime.GOOS == "windows" {
		newPath += ";" + dirPath
	} else {
		newPath += ":" + dirPath
	}

	// 5. Set the new PATH value.
	err := os.Setenv("PATH", newPath)
	if err != nil {
		return fmt.Errorf("error setting PATH: %w", err)
	}
	return nil
}

func main() {
	var input string
	fmt.Print("This program will install the Stream compiler, interpreter, and devkit. You will be able\n" +
		"to type 'stream' in the terminal to open the Stream devkit, and from there will be able to\n" +
		"type 'help' to access any information about the Stream system.\n\n" +
		"Do you wish to proceed? [y/n]")
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Print("Failed to read user input, defaulting to installation.")
		return
	}
	if strings.HasPrefix(strings.ToLower(input), "n") {
		fmt.Print("Installation cancelled.")
		return
	}

	if runtime.GOOS != "windows" {
		fmt.Println("This script is designed for Windows and uses 'choco'.")
		return
	}

	fmt.Println("Generating bin directory")
	mkdir("C:\\Stream")
	mkdir("C:\\Stream\\bin")

	// Set up Python

	fmt.Println("\nAttempting to install Python 3...")
	cmd("choco", []string{"install", "python3", "-y"}, "")

	fmt.Println("\nAttempting to install PyInstaller...")
	cmd("pip", []string{"install", "pyinstaller"}, "")

	// Generate interpreter

	pythonFile := tempPythonFile(pythonCode)
	if pythonFile == "" {
		return
	}
	defer os.Remove(pythonFile)
	fmt.Println("\nTemporary Python script created:", pythonFile)

	fmt.Println("\nAttempting to run PyInstaller...")
	pyInstallerArgs := []string{
		"--onefile",
		pythonFile,
	}
	cmd("pyinstaller", pyInstallerArgs, "")

	fmt.Printf("\nPyInstaller process finished. Executable (if successful) should be in the 'dist' subdirectory.\n")

	err = mv("./dist\\temp_script.exe", "C:\\Stream\\bin\\stream-e.exe")
	if err != nil {
		fmt.Printf("\nError moving executable to \"C:\\Stream\\bin\\stream-e.exe\"")
		return
	}

	fmt.Printf("\nSuccessfully moved executable to \"C:\\Stream\\bin\\stream-e.exe\"")

	// Generate devkit
	pythonFile = tempPythonFile(devkit_code)
	if pythonFile == "" {
		return
	}
	defer os.Remove(devkit_code)
	fmt.Println("\nTemporary Python script created:", pythonFile)

	fmt.Println("\nAttempting to run PyInstaller...")
	pyInstallerArgs = []string{
		"--onefile",
		pythonFile,
	}
	cmd("pyinstaller", pyInstallerArgs, "")

	fmt.Printf("\nPyInstaller process finished. Executable (if successful) should be in the 'dist' subdirectory.\n")

	err = mv("./dist\\temp_script.exe", "C:\\Stream\\bin\\stream.exe")
	if err != nil {
		fmt.Printf("\nError moving executable to \"C:\\Stream\\bin\\stream.exe\"")
		return
	}

	fmt.Printf("\nSuccessfully moved executable to \"C:\\Stream\\bin\\stream.exe\"")

	// Generate web
	pythonFile = tempPythonFile(web_code)
	if pythonFile == "" {
		return
	}
	defer os.Remove(web_code)
	fmt.Println("\nTemporary Python script created:", pythonFile)

	fmt.Println("\nAttempting to run PyInstaller...")
	pyInstallerArgs = []string{
		"--onefile",
		pythonFile,
	}
	cmd("pyinstaller", pyInstallerArgs, "")

	fmt.Printf("\nPyInstaller process finished. Executable (if successful) should be in the 'dist' subdirectory.\n")

	err = mv("./dist\\temp_script.exe", "C:\\Stream\\bin\\stream-web.exe")
	if err != nil {
		fmt.Printf("\nError moving executable to \"C:\\Stream\\bin\\stream-web.exe\"")
		return
	}

	fmt.Printf("\nSuccessfully moved executable to \"C:\\Stream\\bin\\stream-web.exe\"")

	// Generate web compiler
	pythonFile = tempPythonFile(web_compiler_code)
	if pythonFile == "" {
		return
	}
	defer os.Remove(web_compiler_code)
	fmt.Println("\nTemporary Python script created:", pythonFile)

	fmt.Println("\nAttempting to run PyInstaller...")
	pyInstallerArgs = []string{
		"--onefile",
		pythonFile,
	}
	cmd("pyinstaller", pyInstallerArgs, "")

	fmt.Printf("\nPyInstaller process finished. Executable (if successful) should be in the 'dist' subdirectory.\n")

	err = mv("./dist\\temp_script.exe", "C:\\Stream\\bin\\stream-webc.exe")
	if err != nil {
		fmt.Printf("\nError moving executable to \"C:\\Stream\\bin\\stream-webc.exe\"")
		return
	}

	fmt.Printf("\nSuccessfully moved executable to \"C:\\Stream\\bin\\stream-webc.exe\"")

	// Set up Rust
	fmt.Println("Attempting to install Rust (if needed)...")
	cmd("rustup", []string{"default", "stable"}, "")

	// Generate compiler
	tempDir, err := ioutil.TempDir("", "rust_project")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	defer os.RemoveAll(tempDir)

	//  Fill in the Rust code template.
	srcDir := filepath.Join(tempDir, "src") // Create src directory.
	if err := os.Mkdir(srcDir, 0755); err != nil {
		fmt.Println("Error creating src directory:", err)
		return
	}

	err = ioutil.WriteFile(filepath.Join(tempDir, "src", "main.rs"), []byte(rs_code), 0644)
	if err != nil {
		fmt.Println("Error writing Rust source file:", err)
		return
	}

	err = ioutil.WriteFile(filepath.Join(tempDir, "Cargo.toml"), []byte(toml_code), 0644)
	if err != nil {
		fmt.Println("Error writing Toml file:", err)
		return
	}

	fmt.Println("\nAttempting to build Rust project with Cargo...")
	cmd("cargo", []string{"build", "--release"}, tempDir)

	var outFile string
	if runtime.GOOS == "windows" {
		outFile = filepath.Join(tempDir, "target", "release", "Rust-Stream-Compiler.exe")
	} else {
		outFile = filepath.Join(tempDir, "target", "release", "Rust-Stream-Compiler")
	}

	fmt.Printf("\nRust build finished. Executable (if successful) should be at: %s\n", outFile)

	err = mv(outFile, "C:\\Stream\\bin\\stream-fsl.exe")
	if err != nil {
		fmt.Printf("\nError moving executable to \"C:\\Stream\\bin\\stream-fsl.exe\"")
		return
	}

	fmt.Printf("\nSuccessfully moved executable to \"C:\\Stream\\bin\\stream-fsl.exe\"")

	err = addPath("C:\\Stream\\bin")
	if err != nil {
		fmt.Println("\nFailed to add \"C:\\\\Stream\\\\bin\" to the system environment variables")
		return
	}

	fmt.Println("\nSuccessfully added \"C:\\\\Stream\\\\bin\" to the system environment variables")
}
