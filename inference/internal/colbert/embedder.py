import onnxruntime as ort
import numpy as np 
from transformers import AutoTokenizer

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
        
        # pool
        return outputs[0].mean(axis=1).flatten().tolist()




