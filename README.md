
---

## **GoEdit - AI-Powered Terminal File Manager**

GoEdit is a lightweight **terminal-based file editor and AI assistant** written in Go. It enables users to manage files, directories, and system operations using **natural language commands**, while leveraging AI for intelligent file operations.

---

## **✨ Features**

✅ **AI-Powered File Management** – Create, modify, delete, and organize files using natural language.  
✅ **Terminal-Based** – No GUI required; fully functional within the terminal.  
✅ **Full File System Control** – Perform operations on any directory or file.  
✅ **AI Chat Assistant** – Chat directly with the bot in the terminal _(no need to prefix with `ask`)_!  
✅ **Seamless AI Integration** – Uses **Google Gemini API** for intelligent responses.  
✅ **File & Media Analysis** – Analyze images, videos, and document content.  
✅ **Execute System Commands** – Run simple commands directly from the terminal.  
✅ **Cross-Platform Compatibility** – Supports Linux, Windows, and macOS.

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

### **Basic Commands:**
Once the chat starts, you can talk to the AI **without using `ask`**. Just type your command directly:
```sh
Create a new file named notes.txt  
Show all files in the current directory  
Delete the file old_logs.txt  
Analyze the image photo.jpg  
```

### **Advanced Features:**
- **File & Folder Lookup:** Search for specific files and directories.
- **Read & Write Files:** Modify documents using AI assistance.
- **Media Analysis:** Process images and videos.
- **AI Chat:** Ask **any question** that **Google Gemini** can answer!

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

