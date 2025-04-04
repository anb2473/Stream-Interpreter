import pygame


def main(params):
    screen = params['screen']

    pygame.draw.circle(screen, 'red', (100, 100), 100)

    pygame.display.update()

    return 0
