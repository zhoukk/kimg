import React from 'react';
import { Layout } from 'antd';
import styles from './index.css'

const { Header, Content, Footer } = Layout;

function BasicLayout(props) {
  return (
    <Layout className={styles.page}>
      <Header className={styles.header}>
        <div className={styles.title}>Welcome to Kimg</div>
      </Header>
      <Content className={styles.content}>
        {props.children}
      </Content>
      <Footer className={styles.footer}>
        <a href='https://github.com/zhoukk/kimg' target='view_frame'>kimg</a>
      </Footer>
    </Layout>
  );
}

export default BasicLayout;
