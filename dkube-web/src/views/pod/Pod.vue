<template>
    <div class="pod">
        <el-row>
            <el-col :span="24">
                <div>
                    <el-card class="pod-head-card" shadow="never" :body-style="{padding:'10px'}">
                        <el-row>
                            <el-col :span="6">
                                <div>
                                    <span>命名空间: </span>
                                    <el-select v-model="namespaceValue" filterable placeholder="请选择">
                                        <el-option
                                        v-for="(item, index) in namespaceList"
                                        :key="index"
                                        :label="item.metadata.name"
                                        :value="item.metadata.name">
                                        </el-option>
                                    </el-select>
                                </div>
                            </el-col>
                            <el-col :span="2" :offset="16">
                                <div>
                                    <el-button style="border-radius:2px;" icon="Refresh" plain @click="getPods()">刷新</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="pod-head-card" shadow="never" :body-style="{padding:'10px'}">
                        <el-row>
                            <el-col :span="2">
                                <div>
                                    <el-button disabled style="border-radius:2px;" icon="Edit" type="primary">创建</el-button>
                                </div>
                            </el-col>
                            <el-col :span="6">
                                <div>
                                    <el-input class="pod-head-search" clearable placeholder="请输入" v-model="searchInput"></el-input>
                                    <el-button style="border-radius:2px;" icon="Search" type="primary" plain @click="getPods()">搜索</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="pod-body-card" shadow="never" :body-style="{padding:'5px'}">
                        <el-table
                        style="width:100%;font-size:12px;margin-bottom:10px;"
                        :data="podList"
                        v-loading="appLoading"
                        :row-key="getRowKeys"
                        :expand-row-keys="expandKeys"
                        @expand-change="expandChange">
                            <el-table-column width="10"></el-table-column>
                            <el-table-column type="expand">
                                <template #default="props">
                                    <el-tabs v-model="activeName" type="card">
                                        <el-tab-pane label="容器" name="container">
                                            <el-card shadow="never" style="border-radius:1px;" :body-style="{padding:'5px'}">
                                                <el-table
                                                style="width:100%;font-size:12px;"
                                                :data="props.row.spec.containers">
                                                    <el-table-column align=left prop="name" label="容器名"></el-table-column>
                                                    <el-table-column align=left prop="image" label="镜像"></el-table-column>
                                                    <el-table-column align=center label="Pod IP">
                                                        <span>{{ props.row.status.podIP }}</span>
                                                    </el-table-column>
                                                    <el-table-column align=center prop="args" label="启动命令"></el-table-column>
                                                    <el-table-column align=center label="环境变量">
                                                        <template v-slot="scope">
                                                            <el-popover :width="500" placement="left" trigger="hover">
                                                                <el-table style="width:100%;font-size:12px;" size="mini" :show-header="false" :data="scope.row.env">
                                                                    <el-table-column property="name" label="名称"></el-table-column>
                                                                    <el-table-column property="value" label="值"></el-table-column>
                                                                </el-table>
                                                                <template #reference>
                                                                <el-button size="small">此处查看</el-button>
                                                                </template>
                                                            </el-popover>
                                                        </template>
                                                    </el-table-column>
                                                </el-table>
                                            </el-card>
                                        </el-tab-pane>
                                        <el-tab-pane label="日志" name="log">
                                            <el-card shadow="never" style="border-radius:1px;" :body-style="{padding:'5px'}">
                                                <el-row :gutter="10">
                                                    <el-col :span="3">
                                                        <el-select size="small" v-model="containerValue" placeholder="请选择">
                                                            <el-option v-for="item in containerList" :key="item" :value="item">
                                                            </el-option>
                                                        </el-select>
                                                    </el-col>
                                                    <el-col :span="2">
                                                        <el-button style="border-radius:2px;" size="small" type="primary" @click="getPodLog(props.row.metadata.name)">查看</el-button>
                                                    </el-col>
                                                    <el-col :span="24" style="margin-top: 5px">
                                                        <el-card shadow="never" class="pod-body-log-card" :body-style="{padding:'5px'}">
                                                            <span class="pod-body-log-span">{{ logContent }}</span>
                                                        </el-card>
                                                    </el-col>
                                                 </el-row>
                                            </el-card>
                                        </el-tab-pane>
                                        <el-tab-pane label="终端" name="shell">
                                            <el-card shadow="never" style="border-radius:1px;" :body-style="{padding:'5px'}">
                                                <el-row :gutter="10">
                                                    <el-col :span="3">
                                                        <el-select size="small" v-model="containerValue" placeholder="请选择">
                                                            <el-option v-for="item in containerList" :key="item" :value="item">
                                                            </el-option>
                                                        </el-select>
                                                    </el-col>
                                                    <el-col :span="1">
                                                        <el-button style="border-radius:2px;" size="small" type="primary" @click="initSocket(props.row)">连接</el-button>
                                                    </el-col>
                                                    <el-col :span="1">
                                                        <el-button style="border-radius:2px;" size="small" type="danger" @click="closeSocket()">关闭</el-button>
                                                    </el-col>
                                                    <el-col :span="24" style="margin-top: 5px">
                                                        <el-card shadow="never" class="pod-body-shell-card" :body-style="{padding:'5px'}">
                                                            <div id="xterm"></div>
                                                        </el-card>
                                                    </el-col>
                                                 </el-row>
                                            </el-card>
                                        </el-tab-pane>
                                    </el-tabs>
                                </template>
                            </el-table-column>
                            <el-table-column align=left label="Pod名">
                                <template v-slot="scope">
                                    <a class="pod-body-podname" @click="expandMap[scope.row.metadata.name] ? expandChange(scope.row, []) : expandChange(scope.row, [scope.row])">{{ scope.row.metadata.name }}</a>
                                </template>
                            </el-table-column>
                            <el-table-column align=center min-width="150" label="节点">
                                <template v-slot="scope">
                                    <el-tag v-if="scope.row.spec.nodeName !== undefined" type="warning">{{ scope.row.spec.nodeName }}</el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="状态">
                                <template v-slot="scope">
                                    <div :class="{'success-dot':scope.row.status.phase == 'Running', 'warning-dot':scope.row.status.phase == 'Pending', 'error-dot':scope.row.status.phase != 'Running' && scope.row.status.phase != 'Pending'}"></div>
                                    <span :class="{'success-status':scope.row.status.phase == 'Running', 'warning-status':scope.row.status.phase == 'Pending', 'error-status':scope.row.status.phase != 'Running' && scope.row.status.phase != 'Pending'}">{{ scope.row.status.phase }} </span>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="重启数">
                                <template v-slot="scope">
                                    <span>{{ restartTotal(scope) }} </span>
                                </template>
                            </el-table-column>
                            <el-table-column align=center min-width="100" label="创建时间">
                                <template v-slot="scope">
                                    <el-tag type="info">{{ timeTrans(scope.row.metadata.creationTimestamp) }} </el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="操作" width="200">
                                <template v-slot="scope">
                                    <el-button size="small" style="border-radius:2px;" icon="Edit" type="primary" plain @click="getPodDetail(scope)">YAML</el-button>
                                    <el-button size="small" style="border-radius:2px;" icon="Delete" type="danger" @click="handleConfirm(scope, '删除', delPod)">删除</el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                        <el-pagination
                        class="pod-body-pagination"
                        background
                        @size-change="handleSizeChange"
                        @current-change="handleCurrentChange"
                        :current-page="currentPage"
                        :page-sizes="pagesizeList"
                        :page-size="pagesize"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="podTotal">
                        </el-pagination>
                    </el-card>
                </div>
            </el-col>
        </el-row>
        <el-dialog title="YAML信息" v-model="yamlDialog" width="45%" top="5%">
            <codemirror
                :value="contentYaml"
                border
                :options="cmOptions"
                height="500"
                style="font-size:14px;"
                @change="onChange"
            ></codemirror>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="yamlDialog = false">取 消</el-button>
                    <el-button type="primary" @click="updatePod()">更 新</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>

<script>
import common from "../common/Config";
import httpClient from '../../utils/request';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import 'xterm/css/xterm.css';
import 'xterm/lib/xterm.js';
import yaml2obj from 'js-yaml';
import json2yaml from 'json2yaml';
export default {
    data() {
        return {
            cmOptions: common.cmOptions,
            contentYaml: '',
            currentPage: 1,
            pagesize: 10,
            pagesizeList: [10, 20, 30],
            searchInput: '',
            namespaceValue: 'default',
            namespaceList: [],
            namespaceListUrl: common.k8sNamespaceList,
            appLoading: false,
            podList: [],
            podTotal: 0,
            getPodsData: {
                url: common.k8sPodList,
                params: {
                    filter_name: '',
                    namespace: '',
                    page: '',
                    limit: '',
                }
            },
            podDetail: {},
            getPodDetailData: {
                url: common.k8sPodDetail,
                params: {
                    pod_name: '',
                    namespace: ''
                }
            },
            yamlDialog: false,
            updatePodData: {
                url: common.k8sPodUpdate,
                params: {
                    namespace: '',
                    content: ''
                }
            },
            delPodData: {
                url: common.k8sPodDel,
                params: {
                    pod_name: '',
                    namespace: ''
                }
            },
            activeName: 'container',
            expandKeys: [],
            expandMap: {},
            containerList: {},
            containerValue: '',
            getPodContainerData: {
                url: common.k8sPodContainer,
                params: {
                    pod_name: '',
                    namespace: ''
                }
            },
            logContent: '',
            getPodLogData: {
                url: common.k8sPodLog,
                params: {
                    container_name: '',
                    pod_name: '',
                    namespace: ''
                }
            },
            term: null,
            socket: null
        }
    },
    methods: {
        transYaml(content) {
            return json2yaml.stringify(content)
        },
        transObj(content) {
            return yaml2obj.load(content)
        },
        onChange(val) {
            this.contentYaml = val
        },
        handleSizeChange(size) {
            this.pagesize = size;
            this.getPods()
        },
        handleCurrentChange(currentPage) {
            this.currentPage = currentPage;
            this.getPods()
        },
        handleClose(done) {
            this.$confirm('确认关闭？')
            .then(() => {
                done();
            })
            .catch(() => {});
        },
        ellipsis(value) {
            return value.length>15?value.substring(0,15)+'...':value
        },
        timeTrans(timestamp) {
            let date = new Date(new Date(timestamp).getTime() + 8 * 3600 * 1000)
            date = date.toJSON();
            date = date.substring(0, 19).replace('T', ' ')
            return date 
        },
        restartTotal(e) {
            let index, sum = 0
            let containerStatuses = e.row.status.containerStatuses
            for ( index in containerStatuses) {
                sum = sum + containerStatuses[index].restartCount 
            }
            return sum
        },
        getNamespaces() {
            httpClient.get(this.namespaceListUrl)
            .then(res => {
                this.namespaceList = res.data.items
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        getPods() {
            this.appLoading = true
            this.getPodsData.params.filter_name = this.searchInput
            this.getPodsData.params.namespace = this.namespaceValue
            this.getPodsData.params.page = this.currentPage
            this.getPodsData.params.limit = this.pagesize
            httpClient.get(this.getPodsData.url, {params: this.getPodsData.params})
            .then(res => {
                this.podList = res.data.items
                this.podTotal = res.data.total
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.appLoading = false
        },
        getPodDetail(e) {
            this.getPodDetailData.params.pod_name = e.row.metadata.name
            this.getPodDetailData.params.namespace = this.namespaceValue
            httpClient.get(this.getPodDetailData.url, {params: this.getPodDetailData.params})
            .then(res => {
                this.podDetail = res.data
                this.contentYaml = this.transYaml(this.podDetail)
                this.yamlDialog = true
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        updatePod() {
            let content = JSON.stringify(this.transObj(this.contentYaml))
            this.updatePodData.params.namespace = this.namespaceValue
            this.updatePodData.params.content = content
            httpClient.put(this.updatePodData.url, this.updatePodData.params)
            .then(res => {
                this.$message.success({
                message: res.msg
                })
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.yamlDialog = false
        },
        delPod(e) {
            this.delPodData.params.pod_name = e.row.metadata.name
            this.delPodData.params.namespace = this.namespaceValue
            httpClient.delete(this.delPodData.url, {data: this.delPodData.params})
            .then(res => {
                this.getPods()
                this.$message.success({
                message: res.msg
                })
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        handleConfirm(obj, operateName, fn) {
            this.confirmContent = '确认继续 ' + operateName + ' 操作吗？'
            this.$confirm(this.confirmContent,'提示',{
                confirmButtonText: '确定',
                cancelButtonText: '取消',
            })
            .then(() => {
                fn(obj)
                })
            .catch(() => {
                this.$message.info({
                    message: '已取消操作'
                })          
            })
        },
        getRowKeys(row) {
            return row.metadata.name
        },
        expandChange(row, expandedRows) {
            this.expandKeys = []
            this.logContent= ''
            this.containerValue = ''
            this.activeName = 'container'
            if (expandedRows.length > 0) {
                this.expandMap[row.metadata.name] = 1
                this.setExpandMap(row.metadata.name)
                row ? (this.expandKeys.push(row.metadata.name), this. getPodContainer(row)) : ''
            } else {
                this.expandMap[row.metadata.name] = 0
            }
        },
        setExpandMap(podName) {
            let key
            for ( key in this.expandMap ) {
                key !== podName ? this.expandMap[key] = 0 : ''
            }
        },
        getPodContainer(row) {
            this.getPodContainerData.params.pod_name = row.metadata.name
            this.getPodContainerData.params.namespace = this.namespaceValue
            httpClient.get(this.getPodContainerData.url, {params: this.getPodContainerData.params})
            .then(res => {
                this.containerList = res.data
                this.containerValue = this.containerList[0]
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        getPodLog(podName) {
            this.getPodLogData.params.pod_name = podName
            this.getPodLogData.params.container_name = this.containerValue
            this.getPodLogData.params.namespace = this.namespaceValue
            httpClient.get(this.getPodLogData.url, {params: this.getPodLogData.params})
            .then(res => {
                this.logContent = res.data
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        initTerm() {
            this.term = new Terminal({
                rendererType: 'canvas',
                rows: 30,
                cols: 110,
                convertEol: false, 
                scrollback: 10,
                disableStdin: false,
                cursorStyle: 'underline', 
                cursorBlink: true,
                theme: {
                foreground: 'white',
                background: '#060101',
                cursor: 'help'
                }
            });
            this.term.open(document.getElementById('xterm'))
            const fitAddon = new FitAddon()
            this.term.loadAddon(fitAddon)
            fitAddon.fit();
            this.term.focus();
            let _this = this;
            this.term.onData(function (key) {
                let msgOrder = {
                operation: 'stdin',
                data: key,
                };
                _this.socket.send(JSON.stringify(msgOrder));
            });
            let msgOrder2 = {
                operation: 'resize',
                cols: this.term.cols,
                rows: this.term.rows,
            };
            this.socket.send(JSON.stringify(msgOrder2))
        },
        initSocket(row) {
            let terminalWsUrl = common.k8sTerminalWs + "?pod_name=" + row.metadata.name + "&container_name=" + this.containerValue + "&namespace=" + this.namespaceValue
            this.socket = new WebSocket(terminalWsUrl);
            this.socketOnClose();
            this.socketOnOpen();
            this.socketOnMessage();
            this.socketOnError();
        },
        socketOnOpen() {
            this.socket.onopen = () => {
                this.initTerm()
            }
        },
        socketOnMessage() {
            this.socket.onmessage = (msg) => {
                let content = JSON.parse(msg.data)
                this.term.write(content.data)
            }
        },
        socketOnClose() {
            this.socket.onclose = () => {
                this.term.write("链接已关闭")
            }
        },
        socketOnError() {
            this.socket.onerror = () => {
                console.log('socket 链接失败')
            }
        },
        closeSocket() {
            if (this.socket === null) {
                    return 
                }
            this.term.write("链接关闭中。。。")
            this.socket.close()
        }
    },
    watch: {
        namespaceValue: {
            handler() {
                localStorage.setItem('namespace', this.namespaceValue)
                this.currentPage = 1
                this.getPods()
            }
        },
        activeName: {
            handler() {
                if ( this.activeName == 'log' ) {
                    this.expandKeys.length == 1 ? this.getPodLog(this.expandKeys[0]) : ''
                }
            }
        }
    },
    beforeMount() {
        if (localStorage.getItem('namespace') !== undefined && localStorage.getItem('namespace') !== null) {
            this.namespaceValue = localStorage.getItem('namespace')
        }
        this.getNamespaces()
        this.getPods()
    },
    beforeUnmount() {
        if ( this.socket !== null ) {
            this.socket.close()
        }
    },
}
</script>


<style scoped>
    .pod-head-card,.pod-body-card {
        border-radius: 1px;
        margin-bottom: 5px;
    }
    .pod-head-search {
        width:160px;
        margin-right:10px; 
    }
    .pod-body-podname {
        color: #4795EE;
    }
    .pod-body-podname:hover {
        color: rgb(84, 138, 238);
        cursor: pointer;
        font-weight: bold;
    }
    .success-dot{
        display:inline-block;
        width: 7px;
        height:7px;
        background: rgb(27, 202, 21);
        border-radius:50%;
        border:1px solid rgb(27, 202, 21);
        margin-right: 10px;
    }
    .warning-dot{
        display:inline-block;
        width: 7px;
        height:7px;
        background: rgb(233, 200, 16);
        border-radius:50%;
        border:1px solid rgb(233, 200, 16);
        margin-right: 10px;
    }
    .error-dot{
        display:inline-block;
        width: 7px;
        height:7px;
        background: rgb(226, 23, 23);
        border-radius:50%;
        border:1px solid rgb(226, 23, 23);
        margin-right: 10px;
    }
    .success-status {
        color: rgb(27, 202, 21);
    }
    .warning-status {
        color: rgb(233, 200, 16);
    }
    .error-status {
        color: rgb(226, 23, 23);
    }
    :v-deep .el-tabs__item {
        font-size: 12px;
    }
    :v-deep .el-tabs__header {
        margin-bottom: 8px;
    }
    .pod-body-log-card, .pod-body-shell-card {
        border-radius:1px;
        height:600px;
        overflow:auto;
        background-color: #060101;
    }
    .pod-body-log-card {
        color: aliceblue;
    }
    .pod-body-log-span {
        white-space:pre;
    }
</style>