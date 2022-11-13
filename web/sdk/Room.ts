// @ts-nocheck
import { Canvas } from './Canvas';
import { EventCenter } from './EventCenter';
import { FabricObjects } from './Fabric';
import { FabricObject } from './FabricObject';

type optionsType = {
  roomId: string;
  userId: string;
  onlyRead: boolean;
  el: HTMLCanvasElement;
  userName: string;
  token: string;
};
export class Room extends EventCenter {
  public canvasMap: Map<string, Canvas> = new Map();
  public roomId: string;
  public onlyRead: boolean;
  public ws: any = null;
  public currentCanvasId: string = '';
  public userId: string;
  public timeoutObj: any = null;
  public serverTimeoutObj: any = null;
  public reConnectObj: any = null;
  public el: HTMLCanvasElement;
  public userName: string;

  constructor(options: optionsType) {
    super();
    this.roomId = options.roomId;
    this.userId = options.userId;
    this.onlyRead = options.onlyRead || false;
    this.el = options.el;
    this.userName = options.userName;
    this.initWS(options.token);
    this.initBindRoomEvent();
  }

  initWS(token) {
    this.ws = new WebSocket(`ws://81.68.68.216:8888/websocket?token=${token}`);
    this.addHeartBeat();
    this.ws.onmessage = (e: any) => {
      const res = JSON.parse(e.data);
      if (res.modifyOnlyRead === 1) {
        this.emit('modifyOnlyRead');
      } else if (res.contentType === 1) {
        clearTimeout(this.serverTimeoutObj);
      } else if (res.contentType === 2) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const canvas = this.canvasMap.get(res.toWhiteBoard) as Canvas;
          const objects = JSON.parse(res.content);
          canvas.emit('update', { objects: objects });
        }
      } else if (res.contentType === 3) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const canvas = this.canvasMap.get(res.toWhiteBoard) as Canvas;
          const objectData = JSON.parse(res.content);
          const newObject = new FabricObjects[objectData.type](objectData);
          canvas.add(false, newObject);
          newObject.objectId = objectData.objectId;
        }
      } else if (res.contentType === 4) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const canvas = this.canvasMap.get(res.toWhiteBoard) as Canvas;
          for (const object of canvas._objects) {
            if (object.objectId === res.objectId) {
              if (object.active) return;
              const objectData = JSON.parse(res.content) as FabricObject;
              for (const key of Object.keys(objectData)) {
                object[key] = objectData[key];
              }
              canvas.renderAll();
            }
          }
        }
      } else if (res.contentType === 5) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const canvas = this.canvasMap.get(res.toWhiteBoard) as Canvas;
          canvas.delete(res.objectId, false);
        }
      } else if (res.contentType === 6) {
        const objects = JSON.parse(res.content);

        this.emit('switch', { canvasId: res.totoWhiteBoard, objects });
      } else if (res.contentType === 7) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const canvas = this.canvasMap.get(res.toWhiteBoard) as Canvas;
          canvas.emit('updateLock', { objectId: res.objectId, isLock: res.isLock });
        }
      } else if (res.contentType === 8) {
        const canvasId = JSON.parse(res.content).canvasId;
        this.emit('createCanvas', { canvasId });
        this.emit('new', { canvasId, userId: this.userName });
      } else if (res.contentType === 9) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const canvas = this.canvasMap.get(res.toWhiteBoard) as Canvas;
          for (const i of canvas._objects) {
            if (i.objectId === res.objectId) {
              if (!res.isLock) {
                i.emit('canLock');
              } else {
                i.off('canLock');
              }
              break;
            }
          }
        }
      } else if (res.contentType === 10) {
        const leaveUserId = res.userName;
        this.emit('leaveRoom', { leaveUserId });
      } else if (res.contentType === 11) {
        const canvasIds = JSON.parse(res.content).canvasIds;
        for (const i of canvasIds) {
          if (!this.canvasMap.has(i)) {
            this.emit('createCanvas', { canvasId: i });
          }
        }
      } else if (res.contentType === 12) {
        const message = JSON.parse(res.content).message;
        this.emit('customizeMessage', { message });
      } else if (res.contentType === 13) {
        const joinUserId = res.userName;
        this.emit('joinRoom', { joinUserId });
      }
    };
    if (this.ws && this.ws.readyState === 1) {
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          contentType: 11,
        }),
      );
    }
  }

  initBindRoomEvent() {
    this.on('switch', (options: { canvasId: string; objects: string[] }) => {
      if (this.canvasMap.has(options.canvasId)) {
        const temp = this.canvasMap.get(this.currentCanvasId) as Canvas;
        const activeObject = temp.getActiveObject();
        if (activeObject) {
          temp.discardActiveObject();
        }
        temp.clearContext(temp.contextContainer);
        temp.clearContext(temp.contextTop);
        this.currentCanvasId = options.canvasId;
        const canvas = this.canvasMap.get(options.canvasId) as Canvas;
        canvas.emit('update', { objects: options });
      }
    });
    this.on('createCanvas', (param: { canvasId: string }) => {
      const canvas = new Canvas(this.el);
      canvas.canvasId = param.canvasId;
      canvas.onlyRead = this.onlyRead;
      canvas.userId = this.userId;
      this.initBindCanvasEvent(canvas);
      this.canvasMap.set(param.canvasId, canvas);
    });
    this.on('customizeSend', (param: { content: string }) => {
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          contentType: 12,
          content: param.content,
        }),
      );
    });
  }

  sendMessageToAll(message: any) {
    const content = JSON.stringify({ message });
    this.emit('customizeSend', { content });
  }

  addHeartBeat() {
    this.timeoutObj = setInterval(() => {
      if (this.ws && this.ws.readyState === 1) {
        this.ws.send(JSON.stringify({ from: this.userId, contentType: 1 }));
        this.serverTimeoutObj = setTimeout(() => {
          this.closeWs();
        }, 2000);
      }
    }, 3000);
    // setInterval(() => {
    //   if (this.ws && this.ws.readyState === 1) {
    //     this.ws.send(
    //       JSON.stringify({
    //         from: this.userId,
    //         toWhiteBoard: this.currentCanvasId,
    //         contentType: 2,
    //       }),
    //     );
    //     this.ws.send(
    //       JSON.stringify({
    //         from: this.userId,
    //         toRoom: this.roomId,
    //         contentType: 11,
    //       }),
    //     );
    //   }
    // }, 5000);
  }

  closeWs() {
    if (this.ws && this.ws.readyState === 1) {
      this.ws.close();
      clearInterval(this.reConnectObj);
    }
  }

  createCanvas() {
    if (this.ws && this.ws.readyState === 1)
      this.ws.send(JSON.stringify({ from: this.userId, toRoom: this.roomId, contentType: 8 }));
  }

  initBindCanvasEvent(canvas: Canvas) {
    canvas.on('sendLock', (options: { objectId: string }) => {
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          toWhiteBoard: canvas.canvasId,
          objectId: options.objectId,
          isLock: true,
          contentType: 7,
          timestamp: new Date().valueOf(),
        }),
      );
    });
    canvas.on('sendUnLock', (options: { objectId: string }) => {
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          toWhiteBoard: canvas.canvasId,
          objectId: options.objectId,
          isLock: false,
          contentType: 7,
        }),
      );
    });
    canvas.on('updateLock', (options: { objectId: string; isLock: boolean }) => {
      for (const i of canvas._objects) {
        if (i.objectId === options.objectId) {
          i.updateLock(options.isLock);
        }
      }
    });
    canvas.on('object:added', (options: { target: FabricObject }) => {
      options.target.saveState();
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          toWhiteBoard: canvas.canvasId,
          objectId: options.target.objectId,
          content: JSON.stringify(options.target.originalState),
          contentType: 3,
          timestamp: options.target.timestamp,
        }),
      );
    });
    canvas.on('object:delete', (options: { target: FabricObject }) => {
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          toWhiteBoard: canvas.canvasId,
          objectId: options.target.objectId,
          contentType: 5,
          timestamp: options.target.timestamp,
        }),
      );
    });
    canvas.on('object:modified', (options: { target: FabricObject }) => {
      options.target.saveState();
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          toWhiteBoard: canvas.canvasId,
          objectId: options.target.objectId,
          content: JSON.stringify(options.target.originalState),
          contentType: 4,
          timestamp: options.target.timestamp,
        }),
      );
    });
    canvas.on('update', (options: { objects: string[] }) => {
      let objects = options.objects.filter((i) => i).map((i) => JSON.parse(i) as FabricObject);
      let activeObject = canvas.getActiveObject();
      if (activeObject) {
        objects = objects.filter((i) => i.objectId !== activeObject.objectId);
        objects.push(activeObject);
      }
      for (const i of canvas._objects) {
        for (const j of objects) {
          if (i.objectId === j.objectId) {
            for (const key in j) {
              i[key] = j[key];
            }
            break;
          }
        }
      }
      if (this.currentCanvasId === canvas.canvasId) {
        canvas.renderAll();
      }
    });
  }

  switchCanvas(canvasId: string) {
    if (this.canvasMap.has(canvasId)) {
      if (this.currentCanvasId) {
        const temp = this.canvasMap.get(this.currentCanvasId) as Canvas;
        const activeObject = temp.getActiveObject();
        if (activeObject) {
          temp.discardActiveObject();
        }
        temp.clearContext(temp.contextContainer);
        temp.clearContext(temp.contextTop);
      }

      this.currentCanvasId = canvasId;

      const canvas = this.canvasMap.get(canvasId) as Canvas;
      canvas.renderAll();
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          toWhiteBoard: canvasId,
          onlyRead: this.onlyRead,
          contentType: 6,
        }),
      );
    }
  }

  getCanvas(canvasId: string) {
    return this.canvasMap.get(canvasId);
  }

  modifyOnlyRead() {
    if (this.onlyRead) {
      this.canvasMap.forEach((value) => {
        value.onlyRead = false;
      });
    } else {
      this.canvasMap.forEach((value) => {
        value.onlyRead = true;
      });
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          toWhiteBoard: this.currentCanvasId,
          onlyRead: true,
          contentType: 6,
        }),
      );
    }
    this.ws.send(
      JSON.stringify({
        modifyOnlyRead: 1,
      }),
    );
    this.onlyRead = !this.onlyRead;
  }

  kickOutRoom(userId: string) {
    this.ws.send(
      JSON.stringify({
        from: this.userId,
        toRoom: this.roomId,
        userName: userId,
        contentType: 10,
      }),
    );
  }

  exportCanvas(canvasId: string) {
    if (this.canvasMap.has(canvasId)) {
      const temp = this.canvasMap.get(this.currentCanvasId) as Canvas;
      const activeObject = temp.getActiveObject();
      if (activeObject) {
        temp.discardActiveObject();
      }
      const objects = temp._objects.map((i) => {
        i.saveState();
        return i.originalState;
      });
      return btoa(JSON.stringify(objects));
    }
  }

  importCanvas(canvasId: string, canvasData: string) {
    const objects = JSON.parse(atob(canvasData)).map((i) => new FabricObjects[i.type](i));
    console.log('abc', objects);

    if (this.canvasMap.has(canvasId)) {
      const canvas = this.canvasMap.get(canvasId) as Canvas;
      canvas.add(true, ...objects);
    }
  }
}
