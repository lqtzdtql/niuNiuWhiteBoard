import { typeMap } from '@Src/constants/Constants';
import { create, logout, roomList } from '@Src/service/home/HomeService';
import { ICreate, IRoom } from '@Src/service/home/IHomeService';
import { IUserInfo } from '@Src/service/login/ILoginService';
import { Button, Input, InputRef, List, Modal, notification, Select } from 'antd';
import VirtualList from 'rc-virtual-list';
import React, { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './style/Home.less';

const ContainerHeight = 550;

const Home: React.FC = () => {
  const userInfo: IUserInfo = JSON.parse(localStorage.getItem('userInfo') || '');

  const navigate = useNavigate();

  const [data, setData] = useState<IRoom[]>([]);
  const [type, setType] = useState('teaching_room');
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isJoinModalOpen, setIsJoinModalOpen] = useState(false);

  const nameRef = useRef<InputRef>(null);
  const idRef = useRef<InputRef>(null);
  const passwordRef = useRef<InputRef>(null);

  const rooms = [
    { label: '教学房', value: 'teaching_room' },
    { label: '游戏房', value: 'playing_room' },
  ];

  const showCreateModal = () => {
    setIsCreateModalOpen(true);
  };

  const showJoinModal = () => {
    setIsJoinModalOpen(true);
  };

  const handleCreateOk = async () => {
    const name = nameRef.current?.input?.value;
    const params: ICreate = {
      name: name || '',
      type,
    };
    const response = await create(params);
    if (response.code === 200) {
      enterRoom({
        host_name: response.room.host_name,
        uuid: response.room.uuid,
        name: response.room.name,
        type: response.room.type,
      });
    } else {
      notification.open({
        message: '出错了',
        description: response.message,
      });
    }

    setIsCreateModalOpen(false);
  };

  const handleJoinOK = () => {
    const id = idRef.current?.input?.value;
    const password = passwordRef.current?.input?.value || '';

    setIsJoinModalOpen(false);
  };

  const handleCreateCancel = () => {
    setIsCreateModalOpen(false);
  };

  const handleJoinCancel = () => {
    setIsJoinModalOpen(false);
  };

  const appendData = async () => {
    const response = await roomList();
    if (response.code === 200) {
      setData(response.roomlist);
    } else {
      notification.open({
        message: '出错了',
        description: response.message,
      });
    }
    // setData([
    //   { uuid: '123', name: 'abc', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: '123', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'we', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'wfe', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'sc', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'gwe', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'csa', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'wef', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'gwee', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'vs', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'yh', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'rj', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'jyt', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'trhy', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'yth', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: '65j', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'trj', type: 'teaching_room', host_id: 1 },
    //   { uuid: '123', name: 'rtj', type: 'teaching_room', host_id: 1 },
    // ]);
  };

  const enterRoom = (item: IRoom) => {
    navigate(`/room`, { state: [item] });
  };

  const signout = async () => {
    const response = await logout();
    if (response.code === 200) {
      localStorage.removeItem('token');
      localStorage.removeItem('userInfo');
      localStorage.removeItem('remember');
      navigate('/login');
    } else {
      notification.open({
        message: '出错了',
        description: response.message,
      });
    }
  };

  useEffect(() => {
    appendData();
  }, []);

  const onScroll = (e: React.UIEvent<HTMLElement, UIEvent>) => {
    if (e.currentTarget.scrollHeight - e.currentTarget.scrollTop === ContainerHeight) {
      appendData();
    }
  };

  return (
    <div className="main-box">
      <div className="navigator">
        <Button className="logout-button" danger type="primary" onClick={signout}>
          退出登录
        </Button>
        <b>{'用户:' + userInfo.name}</b>
        <Button className="create-button" type="primary" onClick={showCreateModal}>
          创建房间
        </Button>
        <Button className="join-button" type="primary" onClick={showJoinModal}>
          加入房间
        </Button>
        <Button className="refresh-button" type="primary" onClick={appendData}>
          刷新
        </Button>

        <Modal
          centered
          destroyOnClose
          closable={false}
          title="创建房间"
          open={isCreateModalOpen}
          okText="创建"
          cancelText="取消"
          onOk={handleCreateOk}
          onCancel={handleCreateCancel}
        >
          <Input required name="room-name" placeholder="房间名" ref={nameRef} />
          <Select
            className="modal-item"
            defaultValue="教学房"
            style={{ width: 120 }}
            onChange={(value) => {
              setType(value);
            }}
            options={rooms}
          />
        </Modal>

        <Modal
          centered
          destroyOnClose
          title="加入房间"
          open={isJoinModalOpen}
          onOk={handleJoinOK}
          onCancel={handleJoinCancel}
          okText="加入"
          cancelText="取消"
        >
          <Input placeholder="房间号" ref={nameRef} />
          <Input className="modal-item" placeholder="密码(可不填)" />
        </Modal>
      </div>
      <List className="list">
        <VirtualList data={data} height={ContainerHeight} itemHeight={47} itemKey="number" onScroll={onScroll}>
          {(item: IRoom) => (
            <List.Item key={item.name}>
              <List.Item.Meta className="item-title" title={item.name} />
              <List.Item.Meta className="item-title" title={typeMap.get(item.type)} />

              <List.Item.Meta className="item-title" title={'房主：' + item.host_name} />
              <Button type="primary" className="item-content" onClick={() => enterRoom(item)}>
                立即加入
              </Button>
            </List.Item>
          )}
        </VirtualList>
      </List>
    </div>
  );
};

export default Home;
