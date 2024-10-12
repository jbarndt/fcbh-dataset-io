# https://huggingface.co/facebook/mms-1b-all

# conda install datasets
from datasets import load_dataset, Audio

# English
stream_data = load_dataset("mozilla-foundation/common_voice_13_0", "en", split="test", streaming=True, trust_remote_code=True)
stream_data = stream_data.cast_column("audio", Audio(sampling_rate=16000))
en_sample = next(iter(stream_data))["audio"]["array"]
