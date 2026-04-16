chrome.runtime.onInstalled.addListener(() => {
  chrome.storage.local.set({ onboardingCompleted: false });
});

let connectionStatus = 'unknown'; // 'unknown', 'connected', 'disconnected', 'error'
let onboardingCompleted = false;

chrome.storage.local.get(['onboardingCompleted'], (result) => {
  if (result.onboardingCompleted !== undefined) {
    onboardingCompleted = result.onboardingCompleted;
  }
});

chrome.storage.onChanged.addListener((changes) => {
  if (changes.onboardingCompleted) {
    onboardingCompleted = changes.onboardingCompleted.newValue;
  }
});

async function sendNativeMessage(command, payload = {}) {
  return new Promise((resolve) => {
    chrome.runtime.sendNativeMessage('com.adblock.proxy', { command, payload }, (response) => {
      if (chrome.runtime.lastError) {
        console.error("Native message error:", chrome.runtime.lastError.message);
        connectionStatus = 'error';
        resolve(null);
      } else {
        if (response && response.status === 'ok') {
          connectionStatus = 'connected';
        } else {
          connectionStatus = 'error';
        }
        resolve(response);
      }
    });
  });
}

// Periodic ping
setInterval(async () => {
  if (onboardingCompleted) {
    const response = await sendNativeMessage('ping');
    if (response) {
      chrome.action.setIcon({ path: "icons/icon16.svg" });
      chrome.action.setBadgeText({ text: "ON" });
      chrome.action.setBadgeBackgroundColor({ color: "#00FF00" });
    } else {
      chrome.action.setBadgeText({ text: "ERR" });
      chrome.action.setBadgeBackgroundColor({ color: "#FF0000" });
    }
  }
}, 5000);

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === 'getOS') {
    chrome.runtime.getPlatformInfo((info) => {
      sendResponse({ os: info.os }); // 'mac', 'win', 'linux', etc.
    });
    return true;
  }
  if (request.action === 'getOnboardingStatus') {
    sendResponse({ onboardingCompleted, connectionStatus });
    return false;
  }
  if (request.action === 'setOnboardingCompleted') {
    onboardingCompleted = true;
    chrome.storage.local.set({ onboardingCompleted: true }, () => {
      sendResponse({ success: true });
    });
    return true;
  }
  if (request.action === 'pingNative') {
    sendNativeMessage('ping').then((response) => {
      sendResponse({ connected: response !== null });
    });
    return true;
  }
  if (request.action === 'getStatus') {
    sendResponse({ connected: connectionStatus === 'connected' });
    return false;
  }
  if (request.action === 'toggleBlocking') {
    sendNativeMessage('setEnabled', { value: request.enabled }).then((response) => {
      sendResponse({ success: response !== null });
    });
    return true;
  }
  if (request.action === 'updateBlocklist') {
    sendNativeMessage('updateBlocklist').then(() => {
      sendResponse({ success: true });
    });
    return true;
  }
});
