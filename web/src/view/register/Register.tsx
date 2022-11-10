import logo from '@Public/static/img/logo.jpg';
import { IRegister, IRegisterResponse } from '@Src/service/register/IRegisterService';
import { register } from '@Src/service/register/RegisterService';
import { Button, Form, Input, notification, Select } from 'antd';
import React from 'react';
import { useNavigate } from 'react-router-dom';
import './style/Register.less';

const { Option } = Select;

const Register: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();

  const onFinish = async (values: any) => {
    const params: IRegister = {
      mobile: values.prefix + values.mobile,
      name: values.nickname,
      passwd: values.password,
    };
    const response: IRegisterResponse = await register(params);

    if (response.code === 200) {
      navigate('/login', { replace: true });
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
    <div className="register-box">
      <img className="img-logo" src={logo} alt={logo} />
      <Form
        form={form}
        name="register"
        className="register-form"
        onFinish={onFinish}
        initialValues={{
          prefix: '86',
        }}
        scrollToFirstError
      >
        <Form.Item
          name="nickname"
          label="昵称"
          rules={[{ required: true, message: '请输入你的昵称!', whitespace: true }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="mobile"
          label="手机号"
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
          <Input addonBefore={prefixSelector} style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item
          name="password"
          label="密码"
          rules={[
            {
              required: true,
              message: '请输入密码!',
            },
            () => ({
              validator(_, value: string) {
                if (!value || (value.length >= 6 && value.length <= 20)) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error('密码长度在6~20!'));
              },
            }),
          ]}
          hasFeedback
        >
          <Input.Password />
        </Form.Item>

        <Form.Item
          name="confirm"
          label="确认密码"
          dependencies={['password']}
          hasFeedback
          rules={[
            {
              required: true,
              message: '请确认密码!',
            },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue('password') === value) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error('两次密码不一致!'));
              },
            }),
          ]}
        >
          <Input.Password />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" className="register-form-button">
            注册
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
};

export default Register;
