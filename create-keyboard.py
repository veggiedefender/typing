from pathlib import Path
from PIL import Image

base_url = 'https://september-briefly-sorted-sheets.trycloudflare.com/'
image_base_url = f'{base_url}/k'
link_base_url = f'{base_url}/type'

keyboard = [
    ['main'],
    ['q','w','e','r','t','y','u','i','o','p'],
    ['a','s','d','f','g','h','j','k','l'],
    ['shift','z','x','c','v','b','n','m','backspace'],
    ['numerals','comma','space','period','enter']
]

pressable = [
    'q','w','e','r','t','y','u','i','o','p',
    'a','s','d','f','g','h','j','k','l',
    'z','x','c','v','b','n','m','backspace',
    'comma','space','period','enter'
]

image_url_override = {
    'main': f'{base_url}/screen.png'
}

images = Path('./images')
total_width, _ = Image.open(images / 'combined.png').size

html = '<table><tbody><tr><td>'
for row in keyboard:
    total_row = 0
    for key in row:
        key_width, key_height = Image.open(images / f'{key}.png').size
        percent = round(key_width / total_width * 99, 3)
        
        if key in image_url_override:
            src = image_url_override[key]
        else:
            src = f'{image_base_url}/{key}.png'

        img = f'<img src="{src}" width="{percent}%" alt="{key}" align="top">'
        if key in pressable:
            html += f'<a href="{link_base_url}/{key}">{img}</a>'
        else:
            html += f'<a href="#">{img}</a>'
    html += '<br>'

html += '</td></tr></tbody></table>'

print(html)
