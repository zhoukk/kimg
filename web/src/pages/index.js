import React, { PureComponent } from 'react';
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
  wm: false, t: 'kimg', ts: 16, tw: 0, tc: '', tsc: '', tsw: 1, tg: 'se', tx: 10, ty: 10, tr: 0, to: 80,
  r: 0, bc: '', g: false, q: 75, ao: true, st: true, f: 'jpg',
}


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
            min={1}
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
      <FormItem label="Rotate">
        {getFieldDecorator('r', { initialValue: defaultQuery.r })(
          <InputNumber min={0} max={360} />
        )}
      </FormItem>
      <FormItem label="Background" wrapperCol={{ span: 8, offset: 2 }}>
        {getFieldDecorator('bc', { initialValue: defaultQuery.bc })(
          <Input prefix='#' addonAfter={<Icon style={{ color: `#${bc}` }} type="bg-colors" />} />
        )}
      </FormItem>
      <Divider />
      <FormItem label="Gray">
        {getFieldDecorator('g', { initialValue: defaultQuery.g })(
          <Switch checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} defaultChecked={false} />
        )}
      </FormItem>
      <Divider />
      <FormItem label="Quality">
        {getFieldDecorator('q', { initialValue: defaultQuery.q })(
          <IntegerStep />
        )}
      </FormItem>
      <Divider />
      <FormItem label="AutoOrient">
        {getFieldDecorator('ao', { initialValue: defaultQuery.ao })(
          <Switch checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} defaultChecked={true} />
        )}
      </FormItem>
      <FormItem label="Strip">
        {getFieldDecorator('st', { initialValue: defaultQuery.st })(
          <Switch checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} defaultChecked={true} />
        )}
      </FormItem>
      <FormItem label="Format">
        {getFieldDecorator('f', { initialValue: defaultQuery.f })(
          <RadioGroup buttonStyle="solid">
            <RadioButton value='jpg'>JPG</RadioButton>
            <RadioButton value='png'>PNG</RadioButton>
            <RadioButton value='webp'>WEBP</RadioButton>
            <RadioButton value='gif'>GIF</RadioButton>
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
      <FormItem label="Scale Enable">
        {getFieldDecorator('s', { initialValue: defaultQuery.s })(
          <Switch checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} defaultChecked={false} />
        )}
      </FormItem>
      <FormItem label="Scale Mode">
        {getFieldDecorator('sm', { initialValue: defaultQuery.sm })(
          <RadioGroup disabled={!enableScale} buttonStyle="solid">
            <RadioButton value='fit'>FIT</RadioButton>
            <RadioButton value='fill'>FILL</RadioButton>
          </RadioGroup>
        )}
      </FormItem>
      <FormItem label="Scale Width">
        {getFieldDecorator('sw', { initialValue: defaultQuery.sw })(
          <InputNumber min={0} disabled={!enableScale} />
        )}
      </FormItem>
      <FormItem label="Scale Height">
        {getFieldDecorator('sh', { initialValue: defaultQuery.sh })(
          <InputNumber min={0} disabled={!enableScale} />
        )}
      </FormItem>
      <FormItem label="Scale Percent">
        {getFieldDecorator('sp', { initialValue: defaultQuery.sp })(
          <InputNumber min={0} disabled={!enableScale} />
        )}
      </FormItem>
      <FormItem label="Width Percent">
        {getFieldDecorator('swp', { initialValue: defaultQuery.swp })(
          <InputNumber min={0} disabled={!enableScale} />
        )}
      </FormItem>
      <FormItem label="Height Percent">
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
      <FormItem label="Crop Enable">
        {getFieldDecorator('c', { initialValue: defaultQuery.c })(
          <Switch checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} defaultChecked={false} />
        )}
      </FormItem>
      <FormItem label="Crop Gravity">
        {getFieldDecorator('cg', { initialValue: defaultQuery.cg })(
          <RadioGroup buttonStyle="solid" disabled={!enableCrop}>
            <div style={{ lineHeight: 0 }}>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='nw'>NW</RadioButton>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='n'>N</RadioButton>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='ne'>NE</RadioButton>
            </div>
            <div style={{ lineHeight: 0 }}>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='w'>W</RadioButton>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='c'>C</RadioButton>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='e'>E</RadioButton>
            </div>
            <div style={{ lineHeight: 0 }}>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='sw'>SW</RadioButton>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='s'>S</RadioButton>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='se'>SE</RadioButton>
            </div>
          </RadioGroup>
        )}
      </FormItem>
      <FormItem label="Crop Width">
        {getFieldDecorator('cw', { initialValue: defaultQuery.cw })(
          <InputNumber min={0} disabled={!enableCrop} />
        )}
      </FormItem>
      <FormItem label="Crop Height">
        {getFieldDecorator('ch', { initialValue: defaultQuery.ch })(
          <InputNumber min={0} disabled={!enableCrop} />
        )}
      </FormItem>
      <FormItem label="Offset Mode">
        {getFieldDecorator('co', { initialValue: defaultQuery.co })(
          <RadioGroup buttonStyle="solid" disabled={!enableCrop}>
            <div style={{ lineHeight: 0 }}>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='lt'>LT</RadioButton>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='lb'>LB</RadioButton>
            </div>
            <div style={{ lineHeight: 0 }}>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='rt'>RT</RadioButton>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='rb'>RB</RadioButton>
            </div>
          </RadioGroup>
        )}
      </FormItem>
      <FormItem label="Offset X">
        {getFieldDecorator('cx', { initialValue: defaultQuery.cx })(
          <InputNumber min={0} disabled={!enableCrop} />
        )}
      </FormItem>
      <FormItem label="Offset Y">
        {getFieldDecorator('cy', { initialValue: defaultQuery.cy })(
          <InputNumber min={0} disabled={!enableCrop} />
        )}
      </FormItem>
    </Form>
  );
});


const ImageWaterMarkQueryForm = Form.create({
  onValuesChange(props, _, values) {
    const q = {}
    if (values.wm && values.t) {
      q.t = values.t
      if (values.ts > 0) {
        q.ts = values.ts
      }
      if (values.tw > 0) {
        q.tw = values.tw
      }
      if (values.tc.length === 6) {
        q.tc = values.tc
      }
      if (values.tsc.length === 6) {
        q.tsc = values.tsc
      }
      if (values.tsw > 0) {
        q.tsw = values.tsw
      }
      if (values.tg) {
        q.tg = values.tg
      }
      if (values.tx > 0) {
        q.tx = values.tx
      }
      if (values.ty > 0) {
        q.ty = values.ty
      }
      if (values.tr > 0) {
        q.tr = values.tr
      }
      if (values.to > 0) {
        q.to = values.to
      }
    }
    props.onChange('watermark', q)
  },
})(props => {
  const { form } = props
  const { getFieldDecorator, getFieldValue } = form

  const enableWaterMark = getFieldValue('wm')

  let tc = getFieldValue('tc')
  if (!tc || tc.length < 6) {
    tc = '000000'
  }
  let tsc = getFieldValue('tsc')
  if (!tsc || tsc.length < 6) {
    tsc = '000000'
  }

  return (
    <Form labelCol={{ xs: { span: 24 }, sm: { span: 6 } }} wrapperCol={{ xs: { span: 24 }, sm: { span: 16, offset: 2 } }}>
      <FormItem label="WaterMark Enable">
        {getFieldDecorator('wm', { initialValue: defaultQuery.wm })(
          <Switch checkedChildren={<Icon type="check" />} unCheckedChildren={<Icon type="close" />} defaultChecked={false} />
        )}
      </FormItem>
      <FormItem label="Text" wrapperCol={{ span: 8, offset: 2 }}>
        {getFieldDecorator('t', { initialValue: defaultQuery.t })(
          <Input placeholder='text watermark' disabled={!enableWaterMark} />
        )}
      </FormItem>
      <FormItem label="FontSize">
        {getFieldDecorator('ts', { initialValue: defaultQuery.ts })(
          <InputNumber min={0} disabled={!enableWaterMark} />
        )}
      </FormItem>
      <FormItem label="FontWeight">
        {getFieldDecorator('tw', { initialValue: defaultQuery.tw })(
          <InputNumber min={0} disabled={!enableWaterMark} />
        )}
      </FormItem>
      <FormItem label="FontColor" wrapperCol={{ span: 8, offset: 2 }}>
        {getFieldDecorator('tc', { initialValue: defaultQuery.tc })(
          <Input prefix='#' addonAfter={<Icon style={{ color: `#${tc}` }} type="bg-colors" />} disabled={!enableWaterMark} />
        )}
      </FormItem>
      <FormItem label="Stroke Color" wrapperCol={{ span: 8, offset: 2 }}>
        {getFieldDecorator('tsc', { initialValue: defaultQuery.tsc })(
          <Input prefix='#' addonAfter={<Icon style={{ color: `#${tsc}` }} type="bg-colors" />} disabled={!enableWaterMark} />
        )}
      </FormItem>
      <FormItem label="Stroke Width">
        {getFieldDecorator('tsw', { initialValue: defaultQuery.tsw })(
          <InputNumber min={0} disabled={!enableWaterMark} />
        )}
      </FormItem>
      <FormItem label="Gravity">
        {getFieldDecorator('tg', { initialValue: defaultQuery.tg })(
          <RadioGroup buttonStyle="solid" disabled={!enableWaterMark}>
            <div style={{ lineHeight: 0 }}>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='nw'>NW</RadioButton>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='n'>N</RadioButton>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='ne'>NE</RadioButton>
            </div>
            <div style={{ lineHeight: 0 }}>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='w'>W</RadioButton>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='c'>C</RadioButton>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='e'>E</RadioButton>
            </div>
            <div style={{ lineHeight: 0 }}>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='sw'>SW</RadioButton>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='s'>S</RadioButton>
              <RadioButton style={{ width: 50, borderRadius: 0 }} value='se'>SE</RadioButton>
            </div>
          </RadioGroup>
        )}
      </FormItem>
      <FormItem label="X">
        {getFieldDecorator('tx', { initialValue: defaultQuery.tx })(
          <InputNumber min={0} disabled={!enableWaterMark} />
        )}
      </FormItem>
      <FormItem label="Y">
        {getFieldDecorator('ty', { initialValue: defaultQuery.ty })(
          <InputNumber min={0} disabled={!enableWaterMark} />
        )}
      </FormItem>
      <FormItem label="Rotate">
        {getFieldDecorator('tr', { initialValue: defaultQuery.tr })(
          <InputNumber min={0} disabled={!enableWaterMark} />
        )}
      </FormItem>
      <FormItem label="Opacity">
        {getFieldDecorator('to', { initialValue: defaultQuery.to })(
          <InputNumber min={0} disabled={!enableWaterMark} />
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
      message.success(`${info.file.name} uploaded successfully.`);
      this.setState({
        md5sum: info.file.response.md5,
        originUrl: `/image/${info.file.response.md5}?origin=1`,
        originInfo: info.file.response
      })
    } else if (status === 'error') {
      message.error(`${info.file.name} upload failed.`);
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
    fetch(`/image/${md5sum}`, { method: 'DELETE' }).then(() => {
      message.success(`delete image ${md5sum} successfully.`);
      this.setState({
        md5sum: '', originUrl: '', imageUrl: '',
        originInfo: { size: 0, width: 0, height: 0, format: '' },
        imageInfo: { size: 0, width: 0, height: 0, format: '' }
      })
    }).catch(e => {
      message.error(`delete image ${md5sum} failed.`);
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
        <p className="ant-upload-text">Click or drag image file to this area to upload</p>
        <p className="ant-upload-hint">Support for a single upload.</p>
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
            message.error(`load image ${md5sum} failed.`);
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
                    <Statistic value={originInfo.size / 1024} precision={2} title='size' suffix='kb' />
                  </Col>
                  <Col span={6}>
                    <Statistic value={originInfo.width} title='width' />
                  </Col>
                  <Col span={6}>
                    <Statistic value={originInfo.height} title='height' />
                  </Col>
                  <Col span={6}>
                    <Statistic value={originInfo.format} title='format' />
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
                <Button type="primary" shape="round" icon="upload">Upload</Button>
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
              <Panel header="Convert Panel" key="convert" style={{ border: 0 }}>
                <Tabs defaultActiveKey="basic">
                  <TabPane tab="Basic" key="basic">
                    <ImageBasicQueryForm onChange={this.onImageQueryChange} />
                  </TabPane>
                  <TabPane tab="Scale" key="scale">
                    <ImageScaleQueryForm onChange={this.onImageQueryChange} />
                  </TabPane>
                  <TabPane tab="Crop" key="crop">
                    <ImageCropQueryForm onChange={this.onImageQueryChange} />
                  </TabPane>
                  <TabPane tab="WaterMark" key="watermark">
                    <ImageWaterMarkQueryForm onChange={this.onImageQueryChange} />
                  </TabPane>
                </Tabs>
              </Panel>
              <Panel header="Admin Panel" key="admin" style={{ border: 0 }}>
                <Card size='small' title='Load Image From Kimg By Md5sum'>
                  <Search defaultValue={md5sum} placeholder='md5sum' size="large" onSearch={md5sum => {
                    if (md5sum.length !== 32) {
                      message.error(`invalid image md5sum: ${md5sum}.`);
                      return
                    }
                    this.setState({
                      md5sum,
                      originUrl: `/image/${md5sum}?origin=1`,
                    })
                  }} />
                  <Divider>
                    <Popconfirm title="Are you sure delete this image?" onConfirm={this.deleteImage}>
                      <Button icon='delete' disabled={md5sum === ''} >DELETE</Button>
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
                    <Statistic value={imageInfo.size / 1024} precision={2} title='size' suffix='kb' />
                  </Col>
                  <Col span={6}>
                    <Statistic value={imageInfo.width} title='width' />
                  </Col>
                  <Col span={6}>
                    <Statistic value={imageInfo.height} title='height' />
                  </Col>
                  <Col span={6}>
                    <Statistic value={imageInfo.format} title='format' />
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
