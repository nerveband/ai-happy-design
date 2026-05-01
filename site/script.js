document.addEventListener('DOMContentLoaded', () => {

  // =========================================================
  // A) Terminal Typing Animation
  // =========================================================
  const terminalSequences = [
    {
      command: "ahd-figma command text.create '{\"text\":\"Hello World\",\"fontSize\":48}'",
      result: '{"id":"42:1","name":"Hello World","type":"TEXT"}'
    },
    {
      command: 'ahd-figma batch landing-page.json',
      result: '\u2713 42 operations completed'
    },
    {
      command: "ahd-figma command document.accessibility_audit '{\"file\":\"ops.json\"}'",
      result: '\u2713 WCAG checks passed'
    },
    {
      command: "ahd-figma command export.image '{\"nodeId\":\"42:1\",\"scale\":2}'",
      result: '\u2713 Exported to design.png'
    }
  ];

  const terminalOutput = document.getElementById('terminal-output');
  let sequenceIndex = 0;

  function typeCommand(text, callback) {
    let i = 0;

    // Build the prompt line
    terminalOutput.innerHTML = '';
    const promptSpan = document.createElement('span');
    promptSpan.style.color = '#28c840';
    promptSpan.style.fontWeight = '700';
    promptSpan.style.marginRight = '8px';
    promptSpan.textContent = '$';
    terminalOutput.appendChild(promptSpan);

    const commandSpan = document.createElement('span');
    terminalOutput.appendChild(commandSpan);

    const cursor = document.createElement('span');
    cursor.className = 'cursor';
    terminalOutput.appendChild(cursor);

    function typeNext() {
      if (i < text.length) {
        commandSpan.textContent += text[i];
        i++;
        setTimeout(typeNext, 30);
      } else {
        if (callback) callback();
      }
    }

    typeNext();
  }

  function showResult(text, callback) {
    // Remove cursor temporarily
    const cursor = terminalOutput.querySelector('.cursor');
    if (cursor) cursor.remove();

    // Add line break
    terminalOutput.appendChild(document.createElement('br'));

    // Create result element
    const resultSpan = document.createElement('span');
    resultSpan.className = 'terminal-result';
    resultSpan.textContent = text;
    terminalOutput.appendChild(resultSpan);

    // Re-add cursor
    const newCursor = document.createElement('span');
    newCursor.className = 'cursor';
    terminalOutput.appendChild(newCursor);

    // Trigger fade-in on next frame
    requestAnimationFrame(() => {
      resultSpan.classList.add('show');
    });

    if (callback) {
      setTimeout(callback, 2500);
    }
  }

  function runSequence() {
    const seq = terminalSequences[sequenceIndex];
    typeCommand(seq.command, () => {
      setTimeout(() => {
        showResult(seq.result, () => {
          sequenceIndex = (sequenceIndex + 1) % terminalSequences.length;
          runSequence();
        });
      }, 300);
    });
  }

  // Start the terminal animation
  runSequence();


  // =========================================================
  // B) Scroll-Triggered Fade-Up Animations
  // =========================================================
  const fadeObserver = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        entry.target.classList.add('visible');
        fadeObserver.unobserve(entry.target);
      }
    });
  }, {
    threshold: 0.15,
    rootMargin: '0px 0px -50px 0px'
  });

  document.querySelectorAll('.fade-up').forEach((el) => {
    fadeObserver.observe(el);
  });


  // =========================================================
  // C) Speed Bar Animation
  // =========================================================
  const speedSection = document.getElementById('speed');
  if (speedSection) {
    const speedObserver = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          const mpcBar = speedSection.querySelector('.mcp-bar');
          const cliBar = speedSection.querySelector('.cli-bar');
          if (mpcBar) {
            mpcBar.style.width = mpcBar.getAttribute('data-width') + '%';
          }
          if (cliBar) {
            cliBar.style.width = cliBar.getAttribute('data-width') + '%';
          }
          speedObserver.unobserve(entry.target);
        }
      });
    }, {
      threshold: 0.15,
      rootMargin: '0px 0px -50px 0px'
    });

    speedObserver.observe(speedSection);
  }


  // =========================================================
  // D) Tab Switching (DX Section)
  // =========================================================
  const tabButtons = document.querySelectorAll('.tab-btn');
  const tabContents = document.querySelectorAll('.tab-content');

  tabButtons.forEach((btn) => {
    btn.addEventListener('click', () => {
      // Remove active from all buttons
      tabButtons.forEach((b) => b.classList.remove('active'));
      // Add active to clicked button
      btn.classList.add('active');

      // Update ARIA states
      tabButtons.forEach((b) => b.setAttribute('aria-selected', 'false'));
      btn.setAttribute('aria-selected', 'true');

      // Hide all tab content
      tabContents.forEach((tc) => tc.classList.add('hidden'));

      // Show target tab content
      const targetId = btn.getAttribute('data-tab');
      const targetTab = document.getElementById(targetId);
      if (targetTab) {
        targetTab.classList.remove('hidden');
      }
    });
  });


  // =========================================================
  // E) Nav Background on Scroll
  // =========================================================
  const nav = document.getElementById('nav');
  let navTicking = false;

  function updateNav() {
    if (window.scrollY > 50) {
      nav.classList.add('nav-solid');
    } else {
      nav.classList.remove('nav-solid');
    }
    navTicking = false;
  }

  window.addEventListener('scroll', () => {
    if (!navTicking) {
      requestAnimationFrame(updateNav);
      navTicking = true;
    }
  }, { passive: true });

  // Run once on load in case page is already scrolled
  updateNav();

});
