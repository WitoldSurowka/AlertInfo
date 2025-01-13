# AlertInfo

**AlertInfo (AI)** is a plugin that integrates ChatGPT and Google for the **biker_info** system. It allows users to perform web searches and process content using AI.  

**AlertInfo** is also an SMS-based service enabling users to request any information from the internet, delivered at a chosen frequency.

You can try the service yourself by sending an SMS to **(0048) 579 390 848**.

---

## How It Works

### Notification Frequency
At the beginning of your SMS, specify the desired notification frequency. Choose from the following options:

- **`alert1d`**: Receive information every day at 6 PM.  
- **`alerttr`**: Receive information during the workweek.  
- **`alert7d`**: Receive information every Friday at 6 PM.  
- **`alert30d`**: Receive information only on the 28th day of the month.  

### Example Use Cases:
1. **`alerttr Tesla stock price`**  
2. **`alert30d next full moon date`**

---

## The Challenge

Using the web application, OpenAI GPT API does not allow internet searches.  

### The Workaround

By combining the **Google Search API** with the **OpenAI ChatGPT API**, we can achieve a similar result:  

![Integration Diagram](https://github.com/user-attachments/assets/f2eea8de-f393-49b6-8b9b-776144c26421)

---

## API Contract Examples

### **Google Search API**
```bash
curl -X GET "https://www.googleapis.com/customsearch/v1?q=your_query&key=GOOGLE_SEARCH_API_KEY&cx=YOUR_SEARCH_ENGINE_ID"

### **openAI chatGPT API**
```bash
curl "https://api.openai.com/v1/chat/completions" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $OPENAI_API_KEY" \
    -d '{
        "model": "gpt-4o-mini",
        "messages": [
            {
                "role": "system",
                "content": "You are a helpful assistant."
            },
            {
                "role": "user",
                "content": "As an answer, provide me ONLY a link to a web service where I can find an answer to my question: Tesla stock price."
            }
        ]
    }'



