from flask import Flask, request, jsonify
import easyocr
from PIL import Image,ImageEnhance
from io import BytesIO
from pdf2image import convert_from_bytes
import base64
import re

app = Flask(__name__)

# โหลด EasyOCR
reader = easyocr.Reader(['th', 'en'])

@app.route('/ocr', methods=['POST'])
def ocr():
    data = request.get_json()
    if not data or 'base64' not in data:
        return jsonify({'text': "N/A", 'error': "Missing 'base64' key in JSON"})
    
    base64_string = data.get('base64', '')
    if not base64_string:
        return jsonify({'text': "N/A", 'error': "Empty base64 string"})
    
    try:
        data_bytes = base64.b64decode(base64_string)
        if len(data_bytes) == 0:
            return jsonify({'text': "N/A", 'error': "Decoded base64 is empty"})
        
        print(f"Data size: {len(data_bytes)} bytes")
        
        try:
            images = convert_from_bytes(data_bytes, dpi=400)
            img = images[0].convert("RGB")
        except Exception as pdf_error:
            print(f"PDF conversion error: {str(pdf_error)}")  # Log PDF error
            try:
                img = Image.open(BytesIO(data_bytes)).convert("RGB")
            except Exception as img_error:
                return jsonify({'text': "N/A", 'error': f"Failed to process image: {str(img_error)}"})
        
        img.save("test_image.png")
        result = reader.readtext("test_image.png")
        text = " ".join([item[1] for item in result]) if result else "N/A"
        text = re.sub(r'^\d+/\d+\s*', '', text)
        return jsonify({'text': text})
    except Exception as e:
        print(f"General error: {str(e)}") 
        return jsonify({'text': "N/A", 'error': str(e)})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8866)