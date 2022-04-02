import KeyTable from '@novnc/novnc/core/input/keysym';

export enum VNCStatusEnum {
  Disconnected = 'Disconnected',
  Connecting = 'Connecting',
  Connected = 'Connected'
}

export const vncStatusToText = {
  [VNCStatusEnum.Disconnected]: '已断开',
  [VNCStatusEnum.Connecting]: '连接中...',
  [VNCStatusEnum.Connected]: '已连接'
};

export const ShortcutKeyOptions = [
  {
    text: 'Ctrl + Alt + Delete',
    sendKey(rfb) {
      rfb.sendCtrlAltDel();
    }
  },

  {
    text: 'Ctrl + Alt + Backspace',
    sendKey(rfb) {
      rfb.sendKey(KeyTable.XK_Control_L, 'ControlLeft', true);
      rfb.sendKey(KeyTable.XK_Alt_L, 'AltLeft', true);
      rfb.sendKey(KeyTable.XK_BackSpace, 'Backspace', true);

      rfb.sendKey(KeyTable.XK_BackSpace, 'Backspace', false);
      rfb.sendKey(KeyTable.XK_Alt_L, 'AltLeft', false);
      rfb.sendKey(KeyTable.XK_Control_L, 'ControlLeft', false);
    }
  }
];
