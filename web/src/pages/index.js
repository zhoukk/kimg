import React, { PureComponent } from 'react';
import { formatMessage, FormattedMessage } from 'umi-plugin-react/locale';
import { Alert, Button, Card, Collapse, Divider, Empty, Form, Radio, Upload, Popconfirm, Icon, Input, InputNumber, Tabs, Row, Col, Slider, Switch, Statistic, message } from 'antd';

import qs from 'querystring';

const Dragger = Upload.Dragger;
const FormItem = Form.Item;
const RadioGroup = Radio.Group;
const RadioButton = Radio.Button;
const TabPane = Tabs.TabPane;
const Search = Input.Search;
const Panel = Collapse.Panel;

const defaultQuery = {
  s: false, sm: 'fit', sw: 200, sh: 200, sp: 50, swp: 50, shp: 50,
  c: false, cg: 'c', cw: 200, ch: 200, co: 'lt', cx: 10, cy: 10,
  r: 0, bc: '', g: false, q: 75, ao: true, st: true, f: 'jpg',
}

const base64Url = s => Buffer.from(s, 'utf8').toString('base64').replace(/=/g, "").replace(/\+/g, "-").replace(/\//g, "_");

class IntegerStep extends React.Component {
  state = {
    inputValue: this.props.value,
  }

  onChange = value => {
    this.setState({
      inputValue: value,
    });
    this.props.onChange && this.props.onChange(value)
  }

  render() {
    const { inputValue } = this.state;
    return (
      <Row>
        <Col span={12}>
          <Slider
            min={1}
            max={100}
            tipFormatter={value => `${value}%`}
            onChange={this.onChange}
            value={typeof inputValue === 'number' ? inputValue : 0}
          />
        </Col>
        <Col span={4}>
          <InputNumber
            min={0}
            max={100}
            formatter={value => `${value}%`}
            parser={value => value.replace('%', '')}
            style={{ marginLeft: 16 }}
            value={inputValue}
            onChange={this.onChange}
          />
        </Col>
      </Row>
    );
  }
}

const ImageBasicQueryForm = Form.create({
  onValuesChange(props, _, values) {
    const q = {}
    if (values.r > 0) {
      q.r = values.r
    }
    if (values.bc.length === 6) {
      q.bc = values.bc
    }
    if (values.g) {
      q.g = 1
    }
    q.q = values.q
    if (values.ao === false) {
      q.ao = 0
    }
    if (values.st === false) {
      q.st = 0
    }
    q.f = values.f
    props.onChange('basic', q)
  },
})(props => {
  const { form } = props
  const { getFieldDecorator, getFieldValue } = form

  let bc = getFieldValue('bc')
  if (!bc || bc.length < 6) {
    bc = '000000'
  }

  return (
    <Form labelCol={{ xs: { span: 24 }, sm: { span: 6 } }} wrapperCol={{ xs: { span: 24 }, sm: { span: 16, offset: 2 } }}>
      <FormItem label={formatMessage({ id: 'ROTATE' })}>
        {getFieldDecorator('r', { initialValue: defaultQuery.r })(
          <InputNumber min={0} max={360} />
        )}
      </FormItem>
      <FormItem label={formatMessage({ id: 'BACKGROUND' })} wrapperCol={{ span: 8, offset: 2 }}>
        {getFieldDecorator('bc', { initialValue: defaultQuery.bc })(
          <Input prefix='#' addonAfter={<Icon style={{ color: `#${bc}` }} type="bg-colors" />} />
        )}
      </FormItem>
      <Divider />
      <FormItem label={formatMessage({ id: 'GRAY' })}>
        {getFieldDecorator('g', { initialValue: defaultQuery.g })(
          <Switch checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} defaultChecked={false} />
        )}
      </FormItem>
      <Divider />
      <FormItem label={formatMessage({ id: 'QUALITY' })}>
        {getFieldDecorator('q', { initialValue: defaultQuery.q })(
          <IntegerStep />
        )}
      </FormItem>
      <Divider />
      <FormItem label={formatMessage({ id: 'AUTO_ORIENT' })}>
        {getFieldDecorator('ao', { initialValue: defaultQuery.ao })(
          <Switch checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} defaultChecked={true} />
        )}
      </FormItem>
      <FormItem label={formatMessage({ id: 'STRIP' })}>
        {getFieldDecorator('st', { initialValue: defaultQuery.st })(
          <Switch checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} defaultChecked={true} />
        )}
      </FormItem>
      <Divider />
      <FormItem label={formatMessage({ id: 'OUTPUT_FORMAT' })}>
        {getFieldDecorator('f', { initialValue: defaultQuery.f })(
          <RadioGroup buttonStyle="solid">
            <RadioButton value='none'><FormattedMessage id='FORMAT_NONE' /></RadioButton>
            <RadioButton value='jpg'><FormattedMessage id='FORMAT_JPG' /></RadioButton>
            <RadioButton value='png'><FormattedMessage id='FORMAT_PNG' /></RadioButton>
            <RadioButton value='webp'><FormattedMessage id='FORMAT_WEBP' /></RadioButton>
            <RadioButton value='gif'><FormattedMessage id='FORMAT_GIF' /></RadioButton>
          </RadioGroup>
        )}
      </FormItem>
    </Form>
  );
});

const ImageScaleQueryForm = Form.create({
  onValuesChange(props, _, values) {
    const q = {}
    if (values.s) {
      q.s = 1
      if (values.sm) {
        q.sm = values.sm
      }
      if (values.sw > 0) {
        q.sw = values.sw
      }
      if (values.sh > 0) {
        q.sh = values.sh
      }
      if (values.sp > 0) {
        q.sp = values.sp
      }
      if (values.swp > 0) {
        q.swp = values.swp
      }
      if (values.shp > 0) {
        q.shp = values.shp
      }
    }
    props.onChange('scale', q)
  },
})(props => {
  const { form } = props
  const { getFieldDecorator, getFieldValue } = form

  const enableScale = getFieldValue('s')

  return (
    <Form labelCol={{ xs: { span: 24 }, sm: { span: 6 } }} wrapperCol={{ xs: { span: 24 }, sm: { span: 16, offset: 2 } }}>
      <FormItem label={formatMessage({ id: 'SCALE_ENABLE' })}>
        {getFieldDecorator('s', { initialValue: defaultQuery.s })(
          <Switch checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} defaultChecked={false} />
        )}
      </FormItem>
      <FormItem label={formatMessage({ id: 'SCALE_MODE' })}>
        {getFieldDecorator('sm', { initialValue: defaultQuery.sm })(
          <RadioGroup disabled={!enableScale} buttonStyle="solid">
            <RadioButton value='fit'>FIT</RadioButton>
            <RadioButton value='fill'>FILL</RadioButton>
          </RadioGroup>
        )}
      </FormItem>
      <FormItem label={formatMessage({ id: 'SCALE_WIDTH' })}>
        {getFieldDecorator('sw', { initialValue: defaultQuery.sw })(
          <InputNumber min={0} disabled={!enableScale} />
        )}
      </FormItem>
      <FormItem label={formatMessage({ id: 'SCALE_HEIGHT' })}>
        {getFieldDecorator('sh', { initialValue: defaultQuery.sh })(
          <InputNumber min={0} disabled={!enableScale} />
        )}
      </FormItem>
      <FormItem label={formatMessage({ id: 'SCALE_PERCENT' })}>
        {getFieldDecorator('sp', { initialValue: defaultQuery.sp })(
          <InputNumber min={0} disabled={!enableScale} />
        )}
      </FormItem>
      <FormItem label={formatMessage({ id: 'SCALE_WIDTH_PERCENT' })}>
        {getFieldDecorator('swp', { initialValue: defaultQuery.swp })(
          <InputNumber min={0} disabled={!enableScale} />
        )}
      </FormItem>
      <FormItem label={formatMessage({ id: 'SCALE_HEIGHT_PERCENT' })}>
        {getFieldDecorator('shp', { initialValue: defaultQuery.shp })(
          <InputNumber min={0} disabled={!enableScale} />
        )}
      </FormItem>
    </Form>
  );
});

const ImageCropQueryForm = Form.create({
  onValuesChange(props, _, values) {
    const q = {}
    if (values.c) {
      q.c = 1
      if (values.cg) {
        q.cg = values.cg
      }
      if (values.cw > 0) {
        q.cw = values.cw
      }
      if (values.ch > 0) {
        q.ch = values.ch
      }
      if (values.co) {
        q.co = values.co
      }
      if (values.cx > 0) {
        q.cx = values.cx
      }
      if (values.cy > 0) {
        q.cy = values.cy
      }
    }
    props.onChange('crop', q)
  },
})(props => {
  const { form } = props
  const { getFieldDecorator, getFieldValue } = form

  const enableCrop = getFieldValue('c')

  return (
    <Form labelCol={{ xs: { span: 24 }, sm: { span: 6 } }} wrapperCol={{ xs: { span: 24 }, sm: { span: 16, offset: 2 } }}>
      <FormItem label={formatMessage({ id: 'CROP_ENABLE' })}>
        {getFieldDecorator('c', { initialValue: defaultQuery.c })(
          <Switch checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} defaultChecked={false} />
        )}
      </FormItem>
      <FormItem label={formatMessage({ id: 'CROP_GRAVITY' })}>
        {getFieldDecorator('cg', { initialValue: defaultQuery.cg })(
          <RadioGroup buttonStyle="solid" disabled={!enableCrop}>
            <div style={{ lineHeight: 0 }}>
              <RadioButton style={{ width: 60, borderRadius: 0 }} value='nw'><FormattedMessage id='GRAVITY_NW' /></RadioButton>
              <RadioButton style={{ width: 60, borderRadius: 0 }} value='n'><FormattedMessage id='GRAVITY_N' /></RadioButton>
              <RadioButton style={{ width: 60, borderRadius: 0 }} value='ne'><FormattedMessage id='GRAVITY_NE' /></RadioButton>
            </div>
            <div style={{ lineHeight: 0 }}>
              <RadioButton style={{ width: 60, borderRadius: 0 }} value='w'><FormattedMessage id='GRAVITY_W' /></RadioButton>
              <RadioButton style={{ width: 60, borderRadius: 0 }} value='c'><FormattedMessage id='GRAVITY_C' /></RadioButton>
              <RadioButton style={{ width: 60, borderRadius: 0 }} value='e'><FormattedMessage id='GRAVITY_E' /></RadioButton>
            </div>
            <div style={{ lineHeight: 0 }}>
              <RadioButton style={{ width: 60, borderRadius: 0 }} value='sw'><FormattedMessage id='GRAVITY_SW' /></RadioButton>
              <RadioButton style={{ width: 60, borderRadius: 0 }} value='s'><FormattedMessage id='GRAVITY_S' /></RadioButton>
              <RadioButton style={{ width: 60, borderRadius: 0 }} value='se'><FormattedMessage id='GRAVITY_SE' /></RadioButton>
            </div>
          </RadioGroup>
        )}
      </FormItem>
      <FormItem label={formatMessage({ id: 'CROP_WIDTH' })}>
        {getFieldDecorator('cw', { initialValue: defaultQuery.cw })(
          <InputNumber min={0} disabled={!enableCrop} />
        )}
      </FormItem>
      <FormItem label={formatMessage({ id: 'CROP_HEIGHT' })}>
        {getFieldDecorator('ch', { initialValue: defaultQuery.ch })(
          <InputNumber min={0} disabled={!enableCrop} />
        )}
      </FormItem>
      <FormItem label={formatMessage({ id: 'OFFSET_MODE' })}>
        {getFieldDecorator('co', { initialValue: defaultQuery.co })(
          <RadioGroup buttonStyle="solid" disabled={!enableCrop}>
            <div style={{ lineHeight: 0 }}>
              <RadioButton style={{ width: 60, borderRadius: 0 }} value='lt'><FormattedMessage id='CROP_LT' /></RadioButton>
              <RadioButton style={{ width: 60, borderRadius: 0 }} value='rt'><FormattedMessage id='CROP_RT' /></RadioButton>
            </div>
            <div style={{ lineHeight: 0 }}>
              <RadioButton style={{ width: 60, borderRadius: 0 }} value='lb'><FormattedMessage id='CROP_LB' /></RadioButton>
              <RadioButton style={{ width: 60, borderRadius: 0 }} value='rb'><FormattedMessage id='CROP_RB' /></RadioButton>
            </div>
          </RadioGroup>
        )}
      </FormItem>
      <FormItem label={formatMessage({ id: 'OFFSET_X' })}>
        {getFieldDecorator('cx', { initialValue: defaultQuery.cx })(
          <InputNumber min={0} disabled={!enableCrop} />
        )}
      </FormItem>
      <FormItem label={formatMessage({ id: 'OFFSET_Y' })}>
        {getFieldDecorator('cy', { initialValue: defaultQuery.cy })(
          <InputNumber min={0} disabled={!enableCrop} />
        )}
      </FormItem>
    </Form>
  );
});

export default class extends PureComponent {

  state = {
    md5sum: '',
    query: '',
    moduleQuery: {},
    imageUrl: '',
    imageInfo: { size: 0, width: 0, height: 0, format: '' },
    originUrl: '',
    originInfo: { size: 0, width: 0, height: 0, format: '' }
  }

  onChange = info => {
    const status = info.file.status;
    if (status === 'done') {
      message.success(formatMessage({ id: 'SUCCESS_UPLOAD_IMAGE' }, { name: info.file.name }));
      this.setState({
        md5sum: info.file.response.md5,
        originUrl: `/image/${info.file.response.md5}?origin=1`,
        originInfo: info.file.response
      })
    } else if (status === 'error') {
      message.error(formatMessage({ id: 'ERR_UPLOAD_IMAGE' }, { name: info.file.name }));
    }
  }

  onImageQueryChange = (module, values) => {
    const { md5sum, moduleQuery } = this.state

    const q = {}
    moduleQuery[module] = values
    Object.keys(moduleQuery).forEach(k => {
      Object.assign(q, moduleQuery[k])
    })
    this.setState({ moduleQuery })

    if (!md5sum) {
      return
    }

    let imageUrl = ''
    const query = qs.stringify(q)
    if (query === '') {
      imageUrl = `/image/${md5sum}`
    } else {
      imageUrl = `/image/${md5sum}?${query}`
    }
    this.setState({ imageUrl, query })
  }

  deleteImage = () => {
    const { md5sum } = this.state
    if ('' === md5sum) {
      return
    }
    fetch(`/image/${md5sum}`, { method: 'DELETE' }).then(() => {
      message.success(formatMessage({ id: 'SUCCESS_DELETE_IMAGE' }, { md5sum }));
      this.setState({
        md5sum: '', originUrl: '', imageUrl: '',
        originInfo: { size: 0, width: 0, height: 0, format: '' },
        imageInfo: { size: 0, width: 0, height: 0, format: '' }
      })
    }).catch(e => {
      message.error(formatMessage({ id: 'ERR_DELETE_IMAGE' }, { md5sum }));
      console.log(e)
    })
  }

  uploadRender() {
    return (
      <Dragger
        showUploadList={false}
        name='file'
        action='/image'
        onChange={this.onChange}
      >
        <p className="ant-upload-drag-icon">
          <Icon type="inbox" />
        </p>
        <p className="ant-upload-text"><FormattedMessage id='TEXT_UPLOAD' /></p>
        <p className="ant-upload-hint"><FormattedMessage id='HINT_UPLOAD' /></p>
      </Dragger>
    )
  }

  imageRender(url) {
    return (
      <img style={{ width: '100%' }} alt='' src={url}
        onLoad={() => {
          const { md5sum, query, originUrl } = this.state
          if (originUrl === url) {
            this.onImageQueryChange()
            const infoUrl = `/info/${md5sum}?origin=1`
            fetch(infoUrl).then(data => data.json()).then(originInfo => {
              this.setState({ originInfo })
            }).catch(e => {
              this.setState({ md5sum: '', originUrl: '' })
              message.error(`fetch ${infoUrl} failed.`);
              console.log(e)
            })
            return
          }
          let infoUrl = ''
          if (query === '') {
            infoUrl = `/info/${md5sum}`
          } else {
            infoUrl = `/info/${md5sum}?${query}`
          }
          fetch(infoUrl).then(data => data.json()).then(imageInfo => {
            this.setState({ imageInfo })
          }).catch(e => {
            message.error(`fetch ${infoUrl} failed.`);
            console.log(e)
          })
        }}
        onError={() => {
          const { md5sum, originUrl } = this.state
          if (url === originUrl) {
            message.error(formatMessage({ id: 'ERR_LOAD_IMAGE' }, { md5sum }));
            this.setState({ md5sum: '', originUrl: '' })
          }
        }}
      />
    )
  }

  render() {

    const { md5sum, originUrl, imageUrl, originInfo, imageInfo } = this.state

    return (
      <div>
        <Row gutter={8}>
          <Col span={8}>
            <Card size='small'
              title={
                <Row>
                  <Col span={6}>
                    <Statistic value={originInfo.size / 1024} precision={2} title={formatMessage({ id: 'SIZE' })} suffix='kb' />
                  </Col>
                  <Col span={6}>
                    <Statistic value={originInfo.width} title={formatMessage({ id: 'WIDTH' })} />
                  </Col>
                  <Col span={6}>
                    <Statistic value={originInfo.height} title={formatMessage({ id: 'HEIGHT' })} />
                  </Col>
                  <Col span={6}>
                    <Statistic value={originInfo.format} title={formatMessage({ id: 'FORMAT' })} />
                  </Col>
                </Row>
              }>
              {originUrl === '' ? this.uploadRender() : this.imageRender(originUrl)}
            </Card>
            {md5sum !== '' && <Alert style={{ marginTop: 10 }} message={`Md5sum: ${md5sum}`} type="info" />}
            <Divider>
              <Upload
                showUploadList={false}
                name='file'
                action='/image'
                onChange={this.onChange}
              >
                <Button type="primary" shape="round" icon="upload"><FormattedMessage id='BTN_UPLOAD' /></Button>
              </Upload>
            </Divider>
          </Col>
          <Col span={8}>
            <Collapse
              accordion
              bordered={false}
              defaultActiveKey={['convert']}
              expandIcon={({ isActive }) => <Icon type="caret-right" rotate={isActive ? 90 : 0} />}
            >
              <Panel header={formatMessage({ id: 'CONVERT_PANEL' })} key="convert" style={{ border: 0 }}>
                <Tabs defaultActiveKey="basic">
                  <TabPane tab={formatMessage({ id: 'TAB_BASIC' })} key="basic">
                    <ImageBasicQueryForm onChange={this.onImageQueryChange} />
                  </TabPane>
                  <TabPane tab={formatMessage({ id: 'TAB_SCALE' })} key="scale">
                    <ImageScaleQueryForm onChange={this.onImageQueryChange} />
                  </TabPane>
                  <TabPane tab={formatMessage({ id: 'TAB_CROP' })} key="crop">
                    <ImageCropQueryForm onChange={this.onImageQueryChange} />
                  </TabPane>
                </Tabs>
              </Panel>
              <Panel header={formatMessage({ id: 'ADMIN_PANEL' })} key="admin" style={{ border: 0 }}>
                <Card size='small' title={formatMessage({ id: 'LOAD_IMAGE_BY_MD5' })}>
                  <Search defaultValue={md5sum} placeholder='md5sum' size="large" onSearch={md5sum => {
                    if (md5sum.length !== 32) {
                      message.error(formatMessage({ id: 'ERR_INVALID_IMAGE' }, { md5sum }));
                      return
                    }
                    this.setState({
                      md5sum,
                      originUrl: `/image/${md5sum}?origin=1`,
                    })
                  }} />
                  <Divider>
                    <Popconfirm title={formatMessage({ id: 'DELETE_IMAGE_CONFIRM' })} onConfirm={this.deleteImage}>
                      <Button icon='delete' disabled={md5sum === ''} ><FormattedMessage id='BTN_DELETE' /></Button>
                    </Popconfirm>
                  </Divider>
                </Card>
              </Panel>
            </Collapse>
          </Col>
          <Col span={8}>
            <Card
              size='small'
              title={
                <Row>
                  <Col span={6}>
                    <Statistic value={imageInfo.size / 1024} precision={2} title={formatMessage({ id: 'SIZE' })} suffix='kb' />
                  </Col>
                  <Col span={6}>
                    <Statistic value={imageInfo.width} title={formatMessage({ id: 'WIDTH' })} />
                  </Col>
                  <Col span={6}>
                    <Statistic value={imageInfo.height} title={formatMessage({ id: 'HEIGHT' })} />
                  </Col>
                  <Col span={6}>
                    <Statistic value={imageInfo.format} title={formatMessage({ id: 'FORMAT' })} />
                  </Col>
                </Row>
              }>
              {imageUrl !== '' ? this.imageRender(imageUrl) : <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />}
            </Card>
            {imageUrl !== '' && <Alert style={{ marginTop: 10 }} message={<a href={imageUrl} style={{ wordWrap: 'break-word' }} alt='' target='view_frame'>{imageUrl}</a>} type="warning" />}
          </Col>
        </Row>
      </div >
    );
  }
}
