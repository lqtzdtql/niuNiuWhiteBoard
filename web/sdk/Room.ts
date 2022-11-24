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
  public canvasMap: Map<string, FabricObject[]> = new Map();
  public roomId: string = '';
  public onlyRead: boolean = false;
  public ws: any = null;
  public currentCanvasId: string = '';
  public userId: string = '';
  public timeoutObj: any = null;
  public serverTimeoutObj: any = null;
  public canvas: Canvas;
  public userName: string;
  public createCanvasTimeObj: any = null;

  constructor(options: optionsType) {
    super();
    this.roomId = options.roomId;
    this.userId = options.userId;
    this.onlyRead = options.onlyRead || false;
    this.canvas = new Canvas(options.el);
    this.canvas.onlyRead = this.onlyRead;
    this.canvas.userId = this.userId;
    this.userName = options.userName;
    this.initWS(options.token);
    this.initBindRoomEvent();
    this.initBindCanvasEvent();
    this.ishostInit();
  }

  initWS(token) {
    this.ws = new WebSocket(`ws://81.68.68.216:8888/websocket?token=${token}`);
    this.addHeartBeat();
    this.ws.onmessage = (e: any) => {
      const res = JSON.parse(e.data);
      if (res.content === 'modifyOnlyRead') {
        this.onlyRead = !this.onlyRead;
        this.canvas.onlyRead = this.onlyRead;
        this.emit('modifyOnlyRead', {});
      } else if (res.contentType === 1) {
        clearTimeout(this.serverTimeoutObj);
      } else if (res.contentType === 2) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const objects = JSON.parse(res.content);
          const canvasId = res.toWhiteBoard;
          this.emit('update', { canvasId, objects });
        }
      } else if (res.contentType === 3) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const objectData = JSON.parse(res.content);
          const newObject = new FabricObjects[objectData.type](objectData);
          newObject.objectId = objectData.objectId;
          newObject.canvas = this.canvas;
          this.canvasMap.get(res.toWhiteBoard)?.push(newObject);
          if (res.toWhiteBoard === this.currentCanvasId) {
            this.canvas.add(false, newObject);
            console.log('hhhhh', newObject);
          }
        }
      } else if (res.contentType === 4) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const objectData = JSON.parse(res.content) as FabricObject;
          if (res.toWhiteBoard === this.currentCanvasId) {
            let needNew = true;
            for (const object of this.canvas._objects) {
              if (object.objectId === res.objectId) {
                needNew = false;
                if (object.active) return;

                for (const key of Object.keys(objectData)) {
                  object[key] = objectData[key];
                }
                this.canvas.renderAll();
                this.canvasMap.set(res.toWhiteBoard, this.canvas._objects);
                break;
              }
            }
            if (needNew) {
              const newObject = new FabricObjects[objectData.type](objectData);
              newObject.objectId = res.objectId;
              newObject.canvas = this.canvas;
              this.canvas.add(false, newObject);
              this.canvasMap.get(res.toWhiteBoard)?.push(newObject);
            }
          } else {
            let needNew = true;
            for (const object of this.canvasMap.get(res.toWhiteBoard) as FabricObject[]) {
              if (object.objectId === res.objectId) {
                needNew = false;
                for (const key of Object.keys(objectData)) {
                  object[key] = objectData[key];
                }
              }
            }
            if (needNew) {
              const newObject = new FabricObjects[objectData.type](objectData);
              newObject.objectId = res.objectId;
              newObject.canvas = this.canvas;
              this.canvasMap.get(res.toWhiteBoard)?.push(newObject);
            }
          }
        }
      } else if (res.contentType === 5) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          if (res.toWhiteBoard === this.currentCanvasId) {
            this.canvas.delete(res.objectId, false);
          }
          for (let i = 0; i < this.canvasMap.get(res.toWhiteBoard)!.length; i++) {
            if (this.canvasMap.get(res.toWhiteBoard)[i].objectId === res.objectId) {
              this.canvasMap.set(res.toWhiteBoard, this.canvasMap.get(res.toWhiteBoard)!.splice(i, 1));
            }
          }
        }
      } else if (res.contentType === 6) {
        const objects = JSON.parse(res.content);
        this.emit('switch', { canvasId: res.toWhiteBoard, objects });
      } else if (res.contentType === 7) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          this.emit('updateLock', { canvasId: res.toWhiteBoard, objectId: res.objectId, isLock: res.isLock });
        }
      } else if (res.contentType === 8) {
        const canvasId = JSON.parse(res.content).canvasId;
        this.canvasMap.set(canvasId, []);
        this.emit('createCanvas', { canvasId });
        this.emit('new', { canvasId, needSwitch: res.userName === this.userId });
      } else if (res.contentType === 9) {
        if (this.currentCanvasId === res.toWhiteBoard) {
          for (const i of this.canvas._objects) {
            if (i.objectId === res.objectId) {
              if (!res.isLock) {
                i.emit('canLock');
                const map = {
                  Line: 1,
                  Rect: 3,
                  Diamond: 4,
                  Triangle: 5,
                  Round: 6,
                  Arrow: 7,
                  Text: 9,
                  Pen: 10,
                };
                const options = JSON.parse(JSON.stringify(i.originalState));
                options.type = map[options.type];
                this.emit('canModify', { objectId: res.objectId, options });
                console.log('hhh', res.objectId, options);
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
        const canvasIds = JSON.parse(res.content);
        console.log('abc', res.content);

        for (const i of canvasIds) {
          if (!this.canvasMap.has(i)) {
            this.canvasMap.set(i, []);
            console.log('map', this.canvasMap);

            this.emit('new', { canvasId: i });
          }
        }
        console.log('abc', 222222);

        this.emit('isNeedCreateFirst');
      } else if (res.contentType === 12) {
        const message = JSON.parse(res.content).message;
        this.emit('customizeMessage', { message });
      } else if (res.contentType === 13) {
        const joinUserId = res.userName;
        this.emit('joinRoom', { joinUserId });
      } else if (res.contentType === 14) {
        this.emit('switchToHost', { canvasId: res.content });
      }
    };
  }

  initBindRoomEvent() {
    this.on('update', (options: { canvasId: string; objects: string[] }) => {
      let objects = options.objects.filter((i) => i).map((i) => JSON.parse(i));
      console.log('hhh', objects);
      if (options.canvasId === this.currentCanvasId) {
        let activeObject = this.canvas.getActiveObject();
        if (activeObject) {
          objects = objects.filter((i) => i.objectId !== activeObject.objectId);
          objects.push(activeObject);
        }
      }
      if (this.canvasMap.has(options.canvasId)) {
        for (const i of objects) {
          let needNew = true;
          for (const j of this.canvasMap.get(options.canvasId) as FabricObject[]) {
            if (i.objectId === j.objectId) {
              if (i.isInGroup) break;
              for (const key in i) {
                j[key] = i[key];
                needNew = false;
              }
              break;
            }
          }
          if (needNew) {
            const newObject = new FabricObjects[i.type](i);
            newObject.objectId = i.objectId;
            newObject.canvas = this.canvas;
            this.canvasMap.get(options.canvasId)!.push(newObject);
          }
        }
      }
      if (this.currentCanvasId === options.canvasId) {
        this.canvas._objects = this.canvasMap.get(options.canvasId) as FabricObject[];
        this.canvas.renderAll();
      }
    });
    this.on('switch', (options: { canvasId: string; objects: string[] }) => {
      if (this.canvasMap.has(options.canvasId)) {
        this.canvas.modifiedList = [];
        this.canvas.modifiedAgainList = [];
        this.emit('canRevoke', { canRevoke: false });
        this.emit('canRedo', { canRedo: false });
        if (this.currentCanvasId) {
          const activeObject = this.canvas.getActiveObject();
          if (activeObject) {
            this.canvas.discardActiveObject();
          }
          this.canvas.clearContext(this.canvas.contextContainer);
          this.canvas.clearContext(this.canvas.contextTop);
          this.canvas.clearContext(this.canvas.contextCache);
        }
        this.currentCanvasId = options.canvasId;
        this.canvas._objects = this.canvasMap.get(options.canvasId) as FabricObject[];
        this.canvas.renderAll();
        this.emit('update', { canvasId: options.canvasId, objects: options });
        this.emit('switchName', { canvasId: options.canvasId });
      }
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
    this.on('hasModified', (param: { objectId: string; options: {} }) => {
      for (const i of this.canvas._objects) {
        console.log('abcd', 1111111, i);

        if (i.objectId === param.objectId) {
          for (const k in param.options) {
            console.log('abcd', i, 2222222, param.options);
            i[k] = param.options[k];
          }
          i.timestamp = new Date().valueOf();
          this.canvas.emit('object:modified', { target: i });
          i.emit('modified');
          this.canvas.renderAll();
          for (const j of this.canvasMap.get(this.currentCanvasId)) {
            if (j.objectId === param.objectId) {
              for (const k of param.options) {
                j[k] = param.options[k];
              }
            }
          }
          break;
        }
      }
    });
    this.on('createCanvas', (options: { canvasId: string }) => {
      this.switchCanvas(options.canvasId);
    });
    this.on('updateLock', (options: { canvasId: string; objectId: string; isLock: boolean }) => {
      if (this.currentCanvasId === options.canvasId) {
        for (const i of this.canvas._objects) {
          if (i.objectId === options.objectId) {
            i.updateLock(options.isLock);
          }
        }
      }

      for (const i of this.canvasMap.get(options.canvasId) as FabricObject[]) {
        if (i.objectId === options.objectId) {
          i.updateLock(options.isLock);
        }
      }
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
    }
  }

  ishostInit() {
    this.on('isNeedCreateFirst', () => {
      if (Object.keys(this.canvasMap).length === 0) {
        this.createCanvas();
      } else {
        this.on('switchToHost', (options: { canvasId: string }) => {
          this.switchCanvas(options.canvasId);
          this.off('switchToHost');
        });
        this.ws.send(
          JSON.stringify({
            from: this.userId,
            toRoom: this.roomId,
            contentType: 14,
          }),
        );
      }
    });
    const i = setInterval(() => {
      if (this.ws && this.ws.readyState === 1) {
        clearInterval(i);
        this.ws.send(
          JSON.stringify({
            from: this.userId,
            toRoom: this.roomId,
            contentType: 11,
          }),
        );
      }
    }, 200);
  }

  createCanvas() {
    this.createCanvasTimeObj = setInterval(() => {
      if (this.ws && this.ws.readyState === 1) {
        this.ws.send(JSON.stringify({ from: this.userId, toRoom: this.roomId, contentType: 8 }));
        clearInterval(this.createCanvasTimeObj);
      }
    }, 200);
  }

  initBindCanvasEvent() {
    this.canvas.on('sendLock', (options: { objectId: string }) => {
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          toWhiteBoard: this.currentCanvasId,
          objectId: options.objectId,
          isLock: true,
          contentType: 7,
          timestamp: new Date().valueOf(),
        }),
      );
    });
    this.canvas.on('sendUnlock', (options: { objectId: string }) => {
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          toWhiteBoard: this.currentCanvasId,
          objectId: options.objectId,
          isLock: false,
          contentType: 7,
        }),
      );
    });

    this.canvas.on('object:added', (options: { target: FabricObject }) => {
      options.target.saveState();
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          toWhiteBoard: this.currentCanvasId,
          objectId: options.target.objectId,
          content: JSON.stringify(options.target.originalState),
          contentType: 3,
          timestamp: options.target.timestamp,
        }),
      );
    });
    this.canvas.on('object:delete', (options: { target: FabricObject }) => {
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          toWhiteBoard: this.currentCanvasId,
          objectId: options.target.objectId,
          contentType: 5,
          timestamp: options.target.timestamp,
        }),
      );
    });
    this.canvas.on('object:modified', (options: { target: FabricObject }) => {
      options.target.saveState();
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          toWhiteBoard: this.currentCanvasId,
          objectId: options.target.objectId,
          content: JSON.stringify(options.target.originalState),
          contentType: 4,
          timestamp: options.target.timestamp,
        }),
      );
    });
  }

  switchCanvas(canvasId: string) {
    if (this.canvasMap.has(canvasId)) {
      if (this.currentCanvasId) {
        const activeObject = this.canvas.getActiveObject();
        if (activeObject) {
          this.canvas.discardActiveObject();
        }
        this.canvas.clearContext(this.canvas.contextContainer);
        this.canvas.clearContext(this.canvas.contextTop);
        this.canvas.clearContext(this.canvas.contextCache);
      }

      this.currentCanvasId = canvasId;

      this.canvas._objects = this.canvasMap.get(canvasId) as FabricObject[];
      this.canvas.modifiedList = [];
      this.canvas.modifiedAgainList = [];
      this.emit('canRevoke', { canRevoke: false });
      this.emit('canRedo', { canRedo: false });
      this.canvas.renderAll();
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

  modifyOnlyRead() {
    this.onlyRead = !this.onlyRead;
    this.canvas.onlyRead = this.onlyRead;
    this.ws.send(
      JSON.stringify({
        content: 'modifyOnlyRead',
        contentType: '',
        from: this.userId,
        toRoom: this.roomId,
      }),
    );
    if (this.onlyRead) {
      if (this.currentCanvasId) {
        this.switchCanvas(this.currentCanvasId);
      }
    }
    this.emit('modifyOnlyRead');
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
    this.ws.close();
  }

  exportCanvas() {
    if (this.currentCanvasId) {
      const activeObject = this.canvas.getActiveObject();
      if (activeObject) {
        this.canvas.discardActiveObject();
      }
    }
    const objects = this.canvas._objects.map((i) => {
      i.saveState();
      return i.originalState;
    });
    return btoa(JSON.stringify(objects));
  }

  importCanvas(canvasData: string) {
    const objects = JSON.parse(atob(canvasData)).map((i) => {
      const newObject = new FabricObjects[i.type](i);
      newObject.canvas = this.canvas;
      return newObject;
    });
    if (this.currentCanvasId) {
      this.canvas.add(false, true, ...objects);
      console.log('ccc', this.canvas._objects);
      this.canvasMap.set(this.currentCanvasId, this.canvas._objects);
      this.canvas.renderAll();
    }
  }
}
