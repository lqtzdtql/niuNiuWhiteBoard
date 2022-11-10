import { IRoom } from '@Src/service/home/IHomeService';
import { makeAutoObservable, runInAction } from 'mobx';
import { QNLocalAudioTrack, QNLocalTrack, QNMicrophoneAudioTrack, QNRemoteAudioTrack } from 'qnweb-rtc';

export class Stream {
  user_id?: string;

  attendee?: IRoom;

  isLocal: boolean = false;

  tag: string = '';

  audioTrack?: QNMicrophoneAudioTrack | QNLocalAudioTrack | QNRemoteAudioTrack;
  audioMuted: boolean = false;

  constructor() {
    makeAutoObservable(this);
  }

  muteTrack(kind: 'audio' | 'video', muted: boolean) {
    console.log('muteTrack', kind, muted);
    runInAction(() => {
      switch (kind) {
        case 'audio':
          if (this.isLocal && this.audioTrack) {
            const localAudioTrack = this.audioTrack as QNLocalTrack;
            localAudioTrack.setMuted(muted);
            this.audioMuted = muted;
          }
          break;
      }
    });
  }

  release() {
    runInAction(() => {
      const localAudioTrack = this.audioTrack as QNLocalTrack;
      if (localAudioTrack) {
        localAudioTrack.destroy();
      }

      this.audioTrack = undefined;
    });
  }
}
