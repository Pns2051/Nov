document.addEventListener('DOMContentLoaded', () => {
  const onboardingView = document.getElementById('onboarding-view');
  const mainView = document.getElementById('main-view');
  
  // Onboarding elements
  const osCommandEl = document.getElementById('os-command');
  const copyBtn = document.getElementById('copy-btn');
  const checkConnectionBtn = document.getElementById('check-connection-btn');
  const connectionResult = document.getElementById('connection-result');
  const finishSetupBtn = document.getElementById('finish-setup-btn');

  // Main view elements
  const statusEl = document.getElementById('status-text');
  const toggleSwitch = document.getElementById('blocking-toggle');
  const updateBtn = document.getElementById('update-blocklist-btn');

  const githubUser = "Pns2051";

  chrome.runtime.sendMessage({ action: 'getOnboardingStatus' }, (status) => {
    if (status.onboardingCompleted) {
      showMainView();
    } else {
      showOnboardingView();
    }
  });

  function showOnboardingView() {
    onboardingView.classList.add('active');
    mainView.classList.remove('active');

    chrome.runtime.sendMessage({ action: 'getOS' }, (info) => {
      let command = "";
      if (info.os === 'mac' || info.os === 'linux') {
        command = `curl -fsSL https://raw.githubusercontent.com/${githubUser}/Nov/main/install.sh | bash`;
      } else if (info.os === 'win') {
        command = `iwr -useb https://raw.githubusercontent.com/${githubUser}/Nov/main/install.ps1 | iex`;
      } else {
        command = `curl -fsSL https://raw.githubusercontent.com/${githubUser}/Nov/main/install.sh | bash`;
      }
      osCommandEl.textContent = command;
    });

    copyBtn.addEventListener('click', () => {
      navigator.clipboard.writeText(osCommandEl.textContent);
      copyBtn.textContent = 'Copied!';
      setTimeout(() => { copyBtn.textContent = 'Copy'; }, 2000);
    });

    checkConnectionBtn.addEventListener('click', () => {
      connectionResult.textContent = 'Checking...';
      chrome.runtime.sendMessage({ action: 'pingNative' }, (res) => {
        if (res && res.connected) {
          connectionResult.textContent = 'Connected successfully!';
          connectionResult.style.color = 'green';
          finishSetupBtn.disabled = false;
        } else {
          connectionResult.textContent = 'Not connected. Is the proxy running?';
          connectionResult.style.color = 'red';
        }
      });
    });

    finishSetupBtn.addEventListener('click', () => {
      chrome.runtime.sendMessage({ action: 'setOnboardingCompleted' }, () => {
        showMainView();
      });
    });
  }

  function showMainView() {
    mainView.classList.add('active');
    onboardingView.classList.remove('active');

    function updateStatus() {
      chrome.runtime.sendMessage({ action: 'getStatus' }, (res) => {
        if (res && res.connected) {
          statusEl.textContent = 'Connected';
          statusEl.style.color = 'green';
          document.getElementById('status-dot').style.background = 'green';
        } else {
          statusEl.textContent = 'Disconnected';
          statusEl.style.color = 'red';
          document.getElementById('status-dot').style.background = 'red';
        }
      });
    }

    updateStatus();
    setInterval(updateStatus, 3000);

    toggleSwitch.addEventListener('change', (e) => {
      chrome.runtime.sendMessage({ action: 'toggleBlocking', enabled: e.target.checked });
    });

    updateBtn.addEventListener('click', () => {
      chrome.runtime.sendMessage({ action: 'updateBlocklist' }, () => {
        alert('Blocklist update started in the background!');
      });
    });
  }
});
