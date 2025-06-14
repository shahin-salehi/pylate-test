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

    def embed_with_tokens(self, text):
        tokens = self.tokenizer(
            text,
            return_tensors="np",
            padding="max_length",
            truncation=True,
            max_length=128,
            return_attention_mask=True,
            return_offsets_mapping=True,
            return_token_type_ids=True
        )

        token_ids = tokens["input_ids"][0]
        token_texts = self.tokenizer.convert_ids_to_tokens(token_ids)

        onnx_inputs = {
            "input_ids": tokens["input_ids"],
            "attention_mask": tokens["attention_mask"]
        }

        outputs = self.session.run(None, onnx_inputs)
        embeddings = outputs[0][0]

        mask = tokens["attention_mask"][0]
        filtered_embeddings = [embeddings[i] for i in range(len(mask)) if mask[i] == 1]
        filtered_tokens = [token_texts[i] for i in range(len(mask)) if mask[i] == 1]

        return filtered_tokens, filtered_embeddings 


    def match(self, query_embeds, doc_embeds):
        query_embeds = np.stack(query_embeds)  # [Q, 128]
        doc_embeds = np.stack(doc_embeds)      # [D, 128]

        # similarity matrix: [Q, D]
        sim = np.matmul(query_embeds, doc_embeds.T)

        # Get best doc token per query token
        max_per_query = np.max(sim, axis=1)
        best_doc_idx_per_query = np.argmax(sim, axis=1)

        # Optionally: get top-k doc tokens overall
        summed_per_doc = np.max(sim, axis=0)  # each doc token's best match
        top_doc_indices = np.argsort(-summed_per_doc)[:5]

        return top_doc_indices  # list of indices in doc_tokens
