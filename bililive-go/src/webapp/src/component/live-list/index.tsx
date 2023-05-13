import React from "react";
import { Button, Divider, PageHeader, Table, Tag } from 'antd';
import PopDialog from '../pop-dialog/index';
import AddRoomDialog from '../add-room-dialog/index';
import API from '../../utils/api';
import './live-list.css';
import { RouteComponentProps } from "react-router-dom";

const api = new API();

const REFRESH_TIME = 3 * 60 * 1000;

interface Props extends RouteComponentProps {
    refresh?: () => void
}

interface IState {
    list: ItemData[],
    addRoomDialogVisible: boolean,
    window: any
}

interface ItemData {
    key: string,
    name: string,
    room: Room,
    address: string,
    tags: string[],
    listening: boolean
    roomId: string
}

interface Room {
    roomName: string;
    url: string;
}

class LiveList extends React.Component<Props, IState> {
    //子控件
    child!: AddRoomDialog;
    //定时器
    timer!: NodeJS.Timeout;

    runStatus = {
        title: '运行状态',
        key: 'tags',
        dataIndex: 'tags',
        render: (tags: { map: (arg0: (tag: any) => JSX.Element) => React.ReactNode; }) => (
            <span>
                {tags.map(tag => {
                    let color = 'green';
                    if (tag === '已停止') {
                        color = 'grey';
                    }
                    if (tag === '监控中') {
                        color = 'green';
                    }
                    if (tag === '录制中') {
                        color = 'red';
                    }

                    return (
                        <Tag color={color} key={tag}>
                            {tag.toUpperCase()}
                        </Tag>
                    );
                })}
            </span>
        ),
    };

    runAction = {
        title: '操作',
        key: 'action',
        dataIndex: 'listening',
        render: (listening: boolean, data: ItemData) => (
            <span>
                <PopDialog
                    title={listening ? "确定停止监控？" : "确定开启监控？"}
                    onConfirm={(e) => {
                        if (listening) {
                            //停止监控
                            api.stopRecord(data.roomId)
                                .then(rsp => {
                                    api.saveSettingsInBackground();
                                    this.refresh();
                                })
                                .catch(err => {
                                    alert(`停止监控失败:\n${err}`);
                                });
                        } else {
                            //开启监控
                            api.startRecord(data.roomId)
                                .then(rsp => {
                                    api.saveSettingsInBackground();
                                    this.refresh();
                                })
                                .catch(err => {
                                    alert(`开启监控失败:\n${err}`);
                                });
                        }
                    }}>
                    <Button type="link" size="small">{listening ? "停止监控" : "开启监控"}</Button>
                </PopDialog>
                <Divider type="vertical" />
                <PopDialog title="确定删除当前直播间？"
                    onConfirm={(e) => {
                        api.deleteRoom(data.roomId)
                            .then(rsp => {
                                api.saveSettingsInBackground();
                                this.refresh();
                            })
                            .catch(err => {
                                alert(`删除直播间失败:\n${err}`);
                            });
                    }}>
                    <Button type="link" size="small">删除</Button>
                </PopDialog>
                <Divider type="vertical" />
                <Button type="link" size="small" onClick={(e) => {
                    this.props.history.push(`/fileList/${data.address}/${data.name}`);
                }}>文件</Button>
            </span>
        ),
    };

    columns = [
        {
            title: '主播名称',
            dataIndex: 'name',
            key: 'name'
        },
        {
            title: '直播间名称',
            dataIndex: 'room',
            key: 'room',
            render: (room: Room) => <a href={room.url} rel="noopener noreferrer" target="_blank">{room.roomName}</a>
        },
        {
            title: '直播平台',
            dataIndex: 'address',
            key: 'address',
        },
        this.runStatus,
        this.runAction
    ];

    smallColums = [
        {
            title: '主播名称',
            dataIndex: 'name',
            key: 'name',
            render: (name: String, data: ItemData) => <a href={data.room.url} rel="noopener noreferrer" target="_blank">{name}</a>
        },
        this.runStatus,
        this.runAction
    ];

    constructor(props: Props) {
        super(props);
        this.state = {
            list: [],
            addRoomDialogVisible: false,
            window: window
        }
    }

    componentDidMount() {
        //refresh data
        this.requestListData();
        this.timer = setInterval(() => {
            this.requestListData();
        }, REFRESH_TIME);
    }

    componentWillUnmount() {
        //clear refresh timer
        clearInterval(this.timer);
    }

    onRef = (ref: AddRoomDialog) => {
        this.child = ref
    }

    /**
     * 当添加房间按钮点击，弹出Dialog
     */
    onAddRoomClick = () => {
        this.child.showModal()
    }

    /**
     * 保存设置至config文件
     */
    onSettingSave = () => {
        api.saveSettings()
            .then((rsp: any) => {
                if (rsp.err_no === 0) {
                    alert("设置保存成功");
                } else {
                    alert("Server Error!");
                }
            }).catch(err => {
                alert(`Server Error!:\n${err}`);
            })
    }

    /**
     * 刷新页面数据
     */
    refresh = () => {
        this.requestListData();
    }

    /**
     * 加载列表数据
     */
    requestListData() {
        api.getRoomList()
            .then(function (rsp: any) {
                if (rsp.length === 0) {
                    return [];
                }
                return rsp.map((item: any, index: number) => {
                    //判断标签状态
                    let tags;
                    if (item.listening === true) {
                        tags = ['监控中'];
                    } else {
                        tags = ['已停止'];
                    }

                    if (item.recording === true) {
                        tags = ['录制中'];
                    }

                    return {
                        key: index + 1,
                        name: item.host_name,
                        room: {
                            roomName: item.room_name,
                            url: item.live_url
                        },
                        address: item.platform_cn_name,
                        tags,
                        listening: item.listening,
                        roomId: item.id
                    };
                });
            })
            .then((data: ItemData[]) => {
                this.setState({
                    list: data
                });
            })
            .catch(err => {
                alert(`加载列表数据失败:\n${err}`);
            });
    }

    render() {
        return (
            <div>
                <div style={{ backgroundColor: '#F5F5F5', }}>
                    <PageHeader
                        ghost={false}
                        title="直播间列表"
                        subTitle="Room List"
                        extra={[
                            <Button key="2" type="default" onClick={this.onSettingSave}>保存设置</Button>,
                            <Button key="1" type="primary" onClick={this.onAddRoomClick}>
                                添加房间
                            </Button>,
                            <AddRoomDialog key="0" ref={this.onRef} refresh={this.refresh} />
                        ]}>
                    </PageHeader>
                </div>
                <Table className="item-pad" columns={
                    (this.state.window.screen.width > 768) ? this.columns : this.smallColums}
                    dataSource={this.state.list}
                    size={(this.state.window.screen.width > 768) ? "default" : "middle"}
                    pagination={false} />
            </div>
        );
    };
}

export default LiveList;
