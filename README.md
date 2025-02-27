

# **GoEdit - AI-Powered Terminal Assistant & Co-Pilot**

GoEdit is an **AI-driven terminal assistant** that acts as your **smart companion** in the command line. It helps with **file operations, system tasks, AI chat, media analysis, and more**, making your terminal experience more interactive and intelligent.

---

## **✨ Features**

✅ **AI Chat Assistant** – Directly chat with AI, no need for prefixes like `ask`.  
✅ **Smart File Operations** – Create, modify, delete, and organize files seamlessly.  
✅ **Terminal Command Execution** – Automate and simplify complex system tasks.  
✅ **Media & Document Analysis** – Understand images, PDFs, and videos.  
✅ **Directory Scanning** – Lookup files and folders with AI assistance.  
✅ **Google Gemini 2.0 Flash Integration** – AI-powered responses for any query.  
✅ **Lightweight & Efficient** – No GUI, runs directly in the terminal.

---

## **📥 Installation**

### **1️⃣ Download the Release**
Go to the **[Releases](https://github.com/your-repo-link/releases)** page and download the appropriate version for your OS:  
🔹 **Linux** → `mybot-linux`  
🔹 **Windows** → `mybot-windows.exe`  
🔹 **MacOS** → `mybot-macos` _(⚠️ Not tested yet – testers needed!)_

### **2️⃣ Run the Application**
#### **Linux/macOS**
```sh
chmod +x mybot-linux  # (For Linux)
./mybot-linux
```
```sh
chmod +x mybot-macos  # (For macOS)
./mybot-macos
```

#### **Windows**
Simply double-click **`mybot-windows.exe`** or run:
```sh
mybot-windows.exe
```

---

## **⚙️ Initial Setup**
When you run GoEdit for the first time, it will **ask for your API key** (for AI functionality).
- The API key is stored securely on your system.
- The storage path will be displayed – **keep it in mind** in case you need to reset your key.

---

## **🌟 Get a Free Google Gemini API Key!**

GoEdit uses **Gemini 2.0 Flash** for AI-powered operations. You can get a **free API key** from **Google AI Studio**:

1. Go to **[Google AI Studio](https://aistudio.google.com/)**
2. Sign in with your **Google Account**
3. Click on **"API Keys"**
4. Generate a **new API key**
5. Copy the key and paste it when prompted by GoEdit

🔹 **Note:** Free-tier users may have some rate limits on API requests.

---

## **🚀 Usage**

### **No Need for `ask` – Just Chat Freely!**
Once the assistant is started, you can directly chat with it in the terminal:
```sh
Hey, what's the weather like today?  
Find me the latest news on AI research.  
Organize my files in the Documents folder.  
Analyze this image: photo.jpg  
```

### **Smart System Operations:**
- **Run terminal commands** intelligently.
- **Read & write files** seamlessly.
- **Perform file analysis** and categorization.
- **Enhance workflow automation** with AI assistance.

---

## **🛠️ Developer Guide**

If you want to **contribute** or understand the code structure, here’s what you need to know:

### **Project Structure**
- **Main Files:**
  - `main.go` – The entry point of the application.
  - `go.mod` & `go.sum` – Go module dependencies.
- **Tool Files:**
  - Various `.go` files like `scan.go`, `tools.go`, etc., handle different features.
  - The filenames indicate their functionality.
- **Release Folder:**
  - Contains the files used for generating the latest **release builds**.
- **Function Test Directory:**
  - Used for dry-run testing before implementation.
- **Other Files (Images, PDFs, Videos):**
  - Just sample/test files added during development.
  - **Not required** for running the project.

### **Branching & Contributions**
- The repository has **only one branch** – everything is directly in the main directory.
- **Fork** the repo, experiment, and send a **Pull Request (PR)** if you improve something!

---

## **⚠️ Known Issues & Future Plans**
- **MacOS version is untested** (Need feedback from testers).
- Some system commands may cause the app to freeze temporarily.
- **Ongoing development** – Expect new features and bug fixes soon!

---

## **🤝 Contribute**
If you find **bugs** or have **suggestions**, feel free to:  
🔹 Open an **issue** 🚀  
🔹 Create a **pull request (PR)** 💡  
🔹 Fork the repo and experiment! 🎨

---

## **📜 License**
This project is licensed under **MIT**. See the [LICENSE](./LICENSE) file for details.

---

## **💬 Need Help?**
Have a question or suggestion? Open an issue, and let’s improve **GoEdit** together! 🚀

---

