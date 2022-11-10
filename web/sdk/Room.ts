import { FabricObjects } from './Fabric';
import { EventCenter } from './EventCenter';
import { Canvas } from './Canvas';
import { FabricObject } from './FabricObject';

type optionsType = {
  roomId: string;
  userId: string;
  onlyRead: boolean;
  el: HTMLCanvasElement;
  elOptions?: {};
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

  constructor(options: optionsType) {
    super();
    this.roomId = options.roomId;
    this.userId = options.userId;
    this.onlyRead = options.onlyRead || false;
    this.initWS();
    this.initBindEvent();
    this.createCanvas(options.el, options?.elOptions);
  }

  initWS() {
    this.ws = new WebSocket('/websocket');
    this.addHeartBeat();
    this.ws.onmessage((e: any) => {
      const res = JSON.parse(e.data);
      if (res.contentType === 1) {
        clearTimeout(this.serverTimeoutObj);
      } else if (res.contentType === 2) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const canvas = this.canvasMap.get(res.toWhiteBoard) as Canvas;
          let objects = JSON.parse(res.content).objects.map((i: string) => JSON.parse(i)) as FabricObject[];
          if (res.toWhiteBoard === this.currentCanvasId) {
            let activeObject = canvas.getActiveObject();
            if (activeObject) {
              objects = objects.filter((i) => i.objectId !== activeObject.objectId);
              objects.push(activeObject);
            }
            canvas._objects = objects;
            canvas.renderAll();
          } else {
            canvas._objects = objects;
          }
        }
      } else if (res.contentType === 3) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const canvas = this.canvasMap.get(res.toWhiteBoard) as Canvas;
          const newObject = JSON.parse(res.content) as FabricObject;
          canvas.add(false, newObject);
        }
      } else if (res.contentType === 4) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const canvas = this.canvasMap.get(res.toWhiteBoard) as Canvas;
          for (const object of canvas._objects) {
            if (object.objectId === res.objectId) {
              if (object.active) return;
              const objectData = JSON.parse(res.content) as FabricObject;
              for (const key in objectData) {
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
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const canvas = this.canvasMap.get(res.toWhiteBoard) as Canvas;
          const objects = JSON.parse(res.content).objects;
          canvas.emit('update', { objects: objects });
        }
      } else if (res.contentType === 7) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const canvas = this.canvasMap.get(res.toWhiteBoard) as Canvas;
          canvas.emit('updateLock', { objectId: res.objectId, isLock: res.isLock });
        }
      } else if (res.contentType === 8) {
        const canvasId = JSON.parse(res.content).canvasId;
        this.emit('createCanvas', { canvasId });
      } else if (res.contentType === 9) {
        if (this.canvasMap.has(res.toWhiteBoard)) {
          const canvas = this.canvasMap.get(res.toWhiteBoard) as Canvas;
          for (const i of canvas._objects) {
            if (i.objectId === res.objectId) {
              i.emit('canLock');
              break;
            }
          }
        }
      } else if (res.contentType === 10) {
        const leaveRoomId = res.leaveRoomId;
        this.emit('leaveRoom', { leaveRoomId });
      }
    });
  }

  initBindEvent() {}

  addHeartBeat() {
    this.timeoutObj = setInterval(() => {
      if (this.ws && this.ws.readyState === 1) {
        this.ws.send(JSON.stringify({ from: this.userId, contentType: 1 }));
        this.serverTimeoutObj = setTimeout(() => {
          this.closeWs();
        }, 2000);
      }
    }, 5000);
  }

  closeWs() {
    if (this.ws && this.ws.readyState === 1) {
      this.ws.close();
      clearInterval(this.reConnectObj);
    }
  }

  createCanvas(el: HTMLCanvasElement, options?: {}) {
    this.ws.send(JSON.stringify({ from: this.userId, toRoom: this.roomId, contentType: 8 }));
    this.on('createCanvas', (param: { canvasId: string }) => {
      const canvas = new Canvas(el, options);
      canvas.canvasId = param.canvasId;
      canvas.onlyRead = this.onlyRead;
      this.initBindCanvasEvent(canvas);
      this.canvasMap.set(param.canvasId, canvas);
      if (this.currentCanvasId) {
        const temp = this.canvasMap.get(this.currentCanvasId) as Canvas;
        temp.clearContext(temp.contextContainer);
        temp.clearContext(temp.contextTop);
      }
      this.currentCanvasId = param.canvasId;
      this.off('createCanvas');
    });
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
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          toWhiteBoard: canvas.canvasId,
          objectId: options.target.objectId,
          content: JSON.stringify(options.target),
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
      this.ws.send(
        JSON.stringify({
          from: this.userId,
          toRoom: this.roomId,
          toWhiteBoard: canvas.canvasId,
          objectId: options.target.objectId,
          content: JSON.stringify(options.target),
          contentType: 4,
          timestamp: options.target.timestamp,
        }),
      );
    });
  }

  switchCanvas(canvasId: string) {
    if (this.canvasMap.has(canvasId)) {
      const temp = this.canvasMap.get(this.currentCanvasId) as Canvas;
      temp.clearContext(temp.contextContainer);
      temp.clearContext(temp.contextTop);
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
      canvas.on('update', (options: { objects: string[] }) => {
        let objects = options.objects.map((i) => JSON.parse(i) as FabricObject);
        let activeObject = canvas.getActiveObject();
        if (activeObject) {
          objects = objects.filter((i) => i.objectId !== activeObject.objectId);
          objects.push(activeObject);
        }
        canvas._objects = objects;
        canvas.renderAll();
        canvas.off('update');
      });
    }
  }

  getCanvas(canvasId: string) {
    return this.canvasMap.get(canvasId);
  }

  modifyOnlyRead(canvasId: string) {
    if (!this.canvasMap.has(canvasId)) return;
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
          toWhiteBoard: canvasId,
          onlyRead: true,
          contentType: 6,
        }),
      );
    }
    this.onlyRead = !this.onlyRead;
  }

  kickOutRoom(userId: string) {
    this.ws.send(
      JSON.stringify({
        from: this.userId,
        toRoom: this.roomId,
        leaveRoomId: userId,
        contentType: 10,
      }),
    );
  }
}
