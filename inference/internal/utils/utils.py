import numpy as np
class Safe:
    def __init__(self) -> None:
        pass

    def make_json_safe(self, obj):
        if isinstance(obj, np.ndarray):
            return obj.tolist()
        elif isinstance(obj, list):
            return [self.make_json_safe(x) for x in obj]
        elif isinstance(obj, dict):
            return {k: self.make_json_safe(v) for k, v in obj.items()}
        return obj

    
