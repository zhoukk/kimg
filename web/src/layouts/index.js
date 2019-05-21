import React from 'react';
import { Layout, Switch } from 'antd';
import { FormattedMessage, setLocale, getLocale } from 'umi-plugin-react/locale';
import pkg from '../../package.json';
import styles from './index.css'

const { Header, Content, Footer } = Layout;

function BasicLayout(props) {
  return (
    <Layout className={styles.page}>
      <Header className={styles.header}>
        <div className={styles.title}>
          <FormattedMessage id="WELCOME_TO_KIMG" />
          <span className={styles.version}> v{pkg.version}</span>
        </div>
      </Header>
      <Content className={styles.content}>
        {props.children}
      </Content>
      <Footer className={styles.footer}>
        <a href='https://github.com/zhoukk/kimg' target='view_frame'>github: zhoukk/kimg</a>
        <span className={styles.lang}>
          <Switch
            checkedChildren="English"
            unCheckedChildren="中文"
            defaultChecked={getLocale() === 'en-US'}
            onChange={checked => setLocale(checked ? 'en-US' : 'zh-CN')}
          />
        </span>
      </Footer>
    </Layout>
  );
}

export default BasicLayout;
