from io import BytesIO

import pytesseract
from fastapi import FastAPI, File, HTTPException, UploadFile
from pdf2image import convert_from_bytes
from PIL import Image

app = FastAPI(title="OCR Service")


@app.get("/health")
def health():
    return {"status": "ok"}


def _ocr_image(image: Image.Image) -> str:
    # Arabic + English (requires tesseract-ocr-ara language pack)
    return pytesseract.image_to_string(image, lang="ara+eng").strip()


@app.post("/ocr")
async def ocr(file: UploadFile = File(...)):
    data = await file.read()
    if not data:
        raise HTTPException(status_code=400, detail="empty file")

    content_type = (file.content_type or "").lower()
    filename = (file.filename or "").lower()

    try:
        if content_type == "application/pdf" or filename.endswith(".pdf"):
            pages = convert_from_bytes(data, first_page=1, last_page=1)
            if not pages:
                raise HTTPException(status_code=400, detail="could not render pdf")
            text = _ocr_image(pages[0])
        elif content_type.startswith("image/") or filename.endswith(
            (".png", ".jpg", ".jpeg", ".webp", ".tif", ".tiff", ".bmp")
        ):
            image = Image.open(BytesIO(data))
            text = _ocr_image(image)
        else:
            raise HTTPException(
                status_code=415,
                detail="unsupported file type; use image or pdf",
            )
    except HTTPException:
        raise
    except Exception as exc:
        raise HTTPException(status_code=500, detail=f"ocr failed: {exc}") from exc

    return {"text": text}
