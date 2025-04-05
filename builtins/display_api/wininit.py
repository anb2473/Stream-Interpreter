import pygame
pygame.init()
def main(params: dict) -> pygame.display:
    name, width, height, icon, resizable = params['name'], params['width'], params['height'], params['icon'], params['resizable']
    if not resizable:
        pygame.display.set_mode((width, height))
    else:
        pygame.display.set_mode((width, height), pygame.RESIZABLE)
    pygame.display.set_caption(name)
    if icon:
        if type(icon) is str:
            pygame.display.set_icon(pygame.image.load(icon))
        else:
            pygame.display.set_icon(icon)
    screen = pygame.display.get_surface()
    return screen