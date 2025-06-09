import onnxruntime as ort
import logging
import numpy as np
from transformers import AutoTokenizer

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)  # Change to DEBUG for more verbose output
## TODO make safe with try thanks
class Embedder:
    def __init__(self, model_path, tokenizer_name="bert-base-uncased"):
        self.session = ort.InferenceSession(model_path, providers=["CPUExecutionProvider"])
        self.tokenizer = AutoTokenizer.from_pretrained(tokenizer_name)

    def Embed(self, text):
        tokens = self.tokenizer(
                text,
                return_tensors="np",
                padding="max_length",
                truncation=True,
                max_length=128, #indicated by model card
        )

        # prep onnx input
        onnx_inputs = {
                "input_ids": tokens["input_ids"],
                "attention_mask": tokens["attention_mask"]
        }

        outputs = self.session.run(None, onnx_inputs)
        embeddings = outputs[0][0]  # shape: (128, 128)

        # Only keep embeddings where attention_mask == 1
        mask = tokens["attention_mask"][0]
        filtered = [embeddings[i] for i in range(len(mask)) if mask[i] == 1]        

        # Convert (1, seq_len, 128) -> List[np.ndarray(128,)]
        return filtered 
