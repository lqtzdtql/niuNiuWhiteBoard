import React from 'react';
import { IParticipant } from './../../service/home/IHomeService';

interface ContextInterface {
  canvasRef?: any;
  roomRef?: any;
  boards?: any;
  canvasClass?: any;
  partic?: IParticipant[];
  checked?: boolean;
  setChecked?: any;
}
export const Context = React.createContext<ContextInterface>({
  canvasRef: null,
  roomRef: null,
  canvasClass: null,
  boards: [],
  partic: [],
  checked: true,
  setChecked: null,
});
