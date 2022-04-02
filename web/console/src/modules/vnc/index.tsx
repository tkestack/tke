import React, { useEffect, useRef, useState } from 'react';
import RFB from '@novnc/novnc/core/rfb';
import { Text, Dropdown, List } from 'tea-component';
import { encode } from 'js-base64';
import { VNCStatusEnum, vncStatusToText, ShortcutKeyOptions } from './constants';

let rfb;

export const VNCPage = () => {
  const searchParams = new URL(location.href).searchParams;
  const clusterId = searchParams.get('clusterId');
  const namespace = searchParams.get('namespace');
  const name = searchParams.get('name');

  const encodePath = `/apis/platform.tkestack.io/v1/clusters/${clusterId}/proxy?path=/apis/subresources.kubevirt.io/v1/namespaces/${namespace}/virtualmachineinstances/${name}/vnc`;

  const vncUrl = `ws://${location.host}/websocket?clusterName=${clusterId}&encodePath=${encode(encodePath)}`;

  const vncBox = useRef(null);

  const [vncStatus, setVncStatus] = useState(VNCStatusEnum.Connecting);

  useEffect(() => {
    rfb = new RFB(vncBox.current, vncUrl);

    rfb.addEventListener('connect', () => {
      setVncStatus(VNCStatusEnum.Connected);
    });

    rfb.addEventListener('disconnect', e => {
      setVncStatus(VNCStatusEnum.Disconnected);
    });
  }, []);

  return (
    <>
      <div
        style={{
          boxSizing: 'border-box',
          width: '100vw',
          height: 30,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          backgroundColor: '#006eff',
          padding: '0 20px'
        }}
      >
        <Dropdown button={<Text style={{ color: 'white' }}>发送远程命令</Text>}>
          <List type="option">
            {ShortcutKeyOptions.map(({ text, sendKey }) => (
              <List.Item onClick={() => sendKey(rfb)}>{text}</List.Item>
            ))}
          </List>
        </Dropdown>

        <Text style={{ color: 'white' }}>状态：{vncStatusToText[vncStatus]}</Text>

        <Text style={{ color: 'white' }}>虚拟机名称：{name}</Text>
      </div>
      <div style={{ width: '100vw', height: 'calc(100vh - 30px)' }} ref={vncBox}></div>
    </>
  );
};
