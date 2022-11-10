import logo from '@Public/static/img/logo.jpg';
import { Button, Form, Input } from 'antd';
import React from 'react';
import { useNavigate } from 'react-router-dom';
import './style/ForgetPassword.less';

const ForgetPassword: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();

  const onFinish = async (values: any) => {
    // console.log('Received values of form: ', values);
    // const params: IForget = {
    //
    // };
    // const response: IForgetResponse = await forget(params);
    // if (response.code === 200) {
    navigate('/login');
    // } else {
    //   alert('出错了');
    // }
  };

  return (
    <div className="forget-box">
      <img className="img-logo" src={logo} alt={logo} />
      <Form
        form={form}
        name="forget"
        className="forget-form"
        onFinish={onFinish}
        initialValues={{
          prefix: '86',
        }}
        scrollToFirstError
      >
        <Form.Item
          name="password"
          label="密码"
          rules={[
            {
              required: true,
              message: '请输入密码!',
            },
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
          <Button type="primary" htmlType="submit" className="forget-form-button">
            修改
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
};

export default ForgetPassword;
