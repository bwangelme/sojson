class SoJSON {
    constructor() {
        this.initElements();
        this.bindEvents();
        this.currentFunction = 'process';
    }

    initElements() {
        // 获取DOM元素
        this.inputText = document.getElementById('input-text');
        this.outputText = document.getElementById('output-text');
        this.processBtn = document.getElementById('process-btn');
        this.btnText = this.processBtn.querySelector('.btn-text');
        this.loading = this.processBtn.querySelector('.loading');
        this.errorMessage = document.getElementById('error-message');
        this.successMessage = document.getElementById('success-message');
        this.indentSelect = document.getElementById('indent-select');
        this.inputCount = document.getElementById('input-count');
        this.outputCount = document.getElementById('output-count');
        
        // 功能按钮
        this.functionBtns = document.querySelectorAll('.btn-function');
        
        // 操作按钮
        this.clearInputBtn = document.getElementById('clear-input');
        this.pasteBtn = document.getElementById('paste-btn');
        this.copyOutputBtn = document.getElementById('copy-output');
        this.downloadBtn = document.getElementById('download-btn');
    }

    bindEvents() {
        // 处理按钮点击
        this.processBtn.addEventListener('click', () => this.processText());
        
        // 功能选择
        this.functionBtns.forEach(btn => {
            btn.addEventListener('click', (e) => this.selectFunction(e.target.dataset.function));
        });
        
        // 输入框变化
        this.inputText.addEventListener('input', () => this.updateCharCount());
        this.inputText.addEventListener('paste', () => {
            setTimeout(() => this.updateCharCount(), 10);
        });
        
        // 操作按钮
        this.clearInputBtn.addEventListener('click', () => this.clearInput());
        this.pasteBtn.addEventListener('click', () => this.pasteFromClipboard());
        this.copyOutputBtn.addEventListener('click', () => this.copyOutput());
        this.downloadBtn.addEventListener('click', () => this.downloadResult());
        
        // 键盘快捷键
        document.addEventListener('keydown', (e) => this.handleKeyboardShortcuts(e));
        
        // 初始化字符计数
        this.updateCharCount();
    }

    selectFunction(func) {
        this.currentFunction = func;
        
        // 更新按钮状态
        this.functionBtns.forEach(btn => {
            btn.classList.remove('active');
            if (btn.dataset.function === func) {
                btn.classList.add('active');
            }
        });
        
        // 处理按钮始终显示"处理"
        this.btnText.textContent = '处理';
        
        // 清除之前的结果和错误
        this.hideMessages();
    }

    async processText() {
        const inputValue = this.inputText.value.trim();
        
        if (!inputValue) {
            this.showError('请输入要处理的文本');
            return;
        }
        
        this.setLoading(true);
        this.hideMessages();
        
        try {
            const result = await this.callAPI(this.currentFunction, {
                text: inputValue,
                indent: parseInt(this.indentSelect.value)
            });
            
            if (result.success) {
                if (this.currentFunction === 'validate') {
                    this.showSuccess(result.valid ? 'JSON格式正确' : 'JSON格式错误');
                    this.outputText.value = result.valid ? '✅ JSON格式正确' : '❌ JSON格式错误';
                } else {
                    this.outputText.value = result.result;
                    this.showSuccess('处理成功');
                }
            } else {
                this.showError(result.error || '处理失败');
                this.outputText.value = '';
            }
        } catch (error) {
            this.showError('网络请求失败: ' + error.message);
            this.outputText.value = '';
        } finally {
            this.setLoading(false);
            this.updateCharCount();
        }
    }

    async callAPI(endpoint, data) {
        const response = await fetch(`/api/${endpoint}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        return await response.json();
    }

    setLoading(isLoading) {
        if (isLoading) {
            this.btnText.style.display = 'none';
            this.loading.style.display = 'inline-flex';
            this.processBtn.disabled = true;
        } else {
            this.btnText.style.display = 'inline';
            this.loading.style.display = 'none';
            this.processBtn.disabled = false;
        }
    }

    showError(message) {
        this.hideMessages();
        this.errorMessage.textContent = message;
        this.errorMessage.style.display = 'block';
        
        // 3秒后自动隐藏
        setTimeout(() => this.hideMessages(), 3000);
    }

    showSuccess(message) {
        this.hideMessages();
        this.successMessage.textContent = message;
        this.successMessage.style.display = 'block';
        
        // 3秒后自动隐藏
        setTimeout(() => this.hideMessages(), 3000);
    }

    hideMessages() {
        this.errorMessage.style.display = 'none';
        this.successMessage.style.display = 'none';
    }

    updateCharCount() {
        this.inputCount.textContent = this.inputText.value.length.toLocaleString();
        this.outputCount.textContent = this.outputText.value.length.toLocaleString();
    }

    clearInput() {
        this.inputText.value = '';
        this.outputText.value = '';
        this.updateCharCount();
        this.hideMessages();
        this.inputText.focus();
    }

    async pasteFromClipboard() {
        try {
            const text = await navigator.clipboard.readText();
            this.inputText.value = text;
            this.updateCharCount();
            this.inputText.focus();
        } catch (error) {
            this.showError('无法访问剪贴板，请手动粘贴');
        }
    }

    async copyOutput() {
        if (!this.outputText.value) {
            this.showError('没有可复制的内容');
            return;
        }
        
        try {
            await navigator.clipboard.writeText(this.outputText.value);
            this.showSuccess('已复制到剪贴板');
        } catch (error) {
            // 降级方案：选中文本
            this.outputText.select();
            try {
                document.execCommand('copy');
                this.showSuccess('已复制到剪贴板');
            } catch (e) {
                this.showError('复制失败，请手动复制');
            }
        }
    }

    downloadResult() {
        if (!this.outputText.value) {
            this.showError('没有可下载的内容');
            return;
        }
        
        const filename = `sojson_result_${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.json`;
        const blob = new Blob([this.outputText.value], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
        
        this.showSuccess('文件已下载');
    }

    handleKeyboardShortcuts(e) {
        // Ctrl/Cmd + Enter: 处理
        if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
            e.preventDefault();
            this.processText();
        }
        
        // Ctrl/Cmd + L: 清空
        if ((e.ctrlKey || e.metaKey) && e.key === 'l') {
            e.preventDefault();
            this.clearInput();
        }
        
        // Ctrl/Cmd + V: 粘贴（在输入框外）
        if ((e.ctrlKey || e.metaKey) && e.key === 'v' && e.target !== this.inputText) {
            e.preventDefault();
            this.pasteFromClipboard();
        }
        
        // Ctrl/Cmd + C: 复制输出（在输出框外）
        if ((e.ctrlKey || e.metaKey) && e.key === 'c' && e.target !== this.outputText && e.target !== this.inputText) {
            e.preventDefault();
            this.copyOutput();
        }
        
        // Escape: 隐藏消息
        if (e.key === 'Escape') {
            this.hideMessages();
        }
    }
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
    new SoJSON();
});
