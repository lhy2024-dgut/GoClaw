/**
 * Copyright(c) Live2D Inc. All rights reserved.
 *
 * Use of this source code is governed by the Live2D Open Software license
 * that can be found at https://www.live2d.com/eula/live2d-open-software-license-agreement_en.html.
 */

import { LAppDelegate } from './lappdelegate';
import * as LAppDefine from './lappdefine';

(window as any).live2dApp = null;

window.addEventListener(
  'load',
  (): void => {
    if (!LAppDelegate.getInstance().initialize()) {
      return;
    }

    LAppDelegate.getInstance().run();

    const app = LAppDelegate.getInstance();
    (window as any).live2dApp = app;
  },
  { passive: true }
);

window.addEventListener(
  'beforeunload',
  (): void => LAppDelegate.releaseInstance(),
  { passive: true }
);

document.addEventListener('DOMContentLoaded', () => {
  const buttons = document.querySelectorAll('#buttons button');
  buttons.forEach(btn => {
    btn.addEventListener('click', () => {
      const group = btn.getAttribute('data-group');
      const index = parseInt(btn.getAttribute('data-index') || '0', 10);
      
      const subdelegates = LAppDelegate.getInstance().getSubdelegates();
      if (subdelegates && subdelegates[0]) {
        subdelegates[0].playMotion(group, index);
      }
    });
  });

  const textInput = document.getElementById('textInput') as HTMLInputElement;
  const speakBtn = document.getElementById('speakBtn');

  let currentAudio: HTMLAudioElement | null = null;

  const playAudio = () => {
    if (currentAudio) {
      currentAudio.pause();
      currentAudio.currentTime = 0;
    }
    currentAudio = new Audio('/speech.mp3');
    currentAudio.play();
  };

  speakBtn?.addEventListener('click', () => {
    const text = textInput.value.trim();
    if (text) {
      playAudio();
    }
  });

  textInput?.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
      const text = textInput.value.trim();
      if (text) {
        playAudio();
      }
    }
  });
});
