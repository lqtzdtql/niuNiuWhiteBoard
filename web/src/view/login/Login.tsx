import { LockOutlined } from '@ant-design/icons';
import logo from '@Public/static/img/logo.jpg';
import { ILogin, IUserInfo } from '@Src/service/login/ILoginService';
import { getUserInfo, login } from '@Src/service/login/LoginService';
import { Button, Checkbox, Form, Input, notification, Select } from 'antd';
import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './style/Login.less';

const { Option } = Select;

const Login = () => {
  const [isLogin, setIsLogin] = useState(false);

  const navigate = useNavigate();

  useEffect(() => {
    if (localStorage.getItem('remember') === 'true') {
      getInfo();
    }
  }, []);

  useEffect(() => {
    if (isLogin) {
      console.log('登录成功,正在跳转界面');
      navigate('/home');
    } else {
      console.log('等待用户登录');
    }
  }, [isLogin]);

  const getInfo = async () => {
    const token = localStorage.getItem('token');
    const userInfo: IUserInfo = JSON.parse(localStorage.getItem('userInfo') || '');
    console.log(userInfo);

    if (token && token !== '' && userInfo && userInfo.uuid != '') {
      const response = await getUserInfo(userInfo.uuid);
      if (response.uuid) {
        setIsLogin(true);
      }
    }
  };

  //提交
  const onFinish = async (values: any) => {
    const mobile = values.prefix + values.mobile;
    const password = values.password;
    const params: ILogin = {
      mobile: mobile,
      passwd: password,
    };
    const response = await login(params);
    if (response.code === 200) {
      if (values.remember) {
        localStorage.setItem('remember', 'true');
      } else {
        localStorage.removeItem('remember');
      }
      localStorage.setItem('userInfo', JSON.stringify(response.user_info));
      setIsLogin(true);
    } else {
      notification.open({
        message: '出错了',
        description: response.message,
      });
    }
  };

  const prefixSelector = (
    <Form.Item name="prefix" noStyle>
      <Select style={{ width: 70 }}>
        <Option value="86">+86</Option>
        <Option value="87">+87</Option>
      </Select>
    </Form.Item>
  );

  return (
    <div className="login-box">
      <img className="img-logo" src={logo} alt={logo} />
      <Form
        name="normal_login"
        className="login-form"
        initialValues={{
          remember: true,
          prefix: '86',
        }}
        onFinish={onFinish}
      >
        <Form.Item
          name="mobile"
          rules={[
            { required: true, message: '请输入手机号!' },
            () => ({
              validator(_, value: string) {
                if (!value || value.length === 11) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error('请输入正确手机号!'));
              },
            }),
          ]}
        >
          <Input addonBefore={prefixSelector} placeholder="手机号" style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item name="password" rules={[{ required: true, message: '请输入密码!' }]}>
          <Input prefix={<LockOutlined className="site-form-item-icon" />} type="password" placeholder="密码" />
        </Form.Item>
        <Form.Item name="remember" valuePropName="checked" noStyle>
          <Checkbox>记住我</Checkbox>
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" className="login-form-button">
            登录
          </Button>
        </Form.Item>
        <Form.Item>
          <Link to="/register">立即注册!</Link>

          {/* <Link className="login-form-forgot" to="/forget">
            忘记密码
          </Link> */}
        </Form.Item>
      </Form>
    </div>
  );
};
export default Login;
