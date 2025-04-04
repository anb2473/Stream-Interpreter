import queue
import threading

import pygame

render_queue = queue.Queue()


def main(params: dict):
    name, width, height = params['name'], params['width'], params['height']
    pygame.display.set_mode((width, height))

    pygame.display.set_caption(name)

    screen = pygame.display.get_surface()

    return screen
