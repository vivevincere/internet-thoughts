To be able to call Google NLP API locally:

1. Run the following locally

```python
import os 
nlp_key_path = os.path.join(os.path.abspath("."), "gkey.json")
print(nlp_key_path) # use this to get the absolute path on your local machine
```

2. Then run
```bash
export GOOGLE_APPLICATION_CREDENTIALS={ PATH }
# where PATH is the output from above
```