import pygame
pygame.init()
def main(params) -> int:
    screen, x, y, radius = params['screen'], params['x'], params['y'], params['radius']
    pygame.draw.circle(screen, 'red', (x, y), radius)
    return 0