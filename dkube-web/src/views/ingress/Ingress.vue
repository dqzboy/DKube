<template>
    <div class="ingress">
        <el-row>
            <el-col :span="24">
                <div>
                    <el-card class="ingress-head-card" shadow="never" :body-style="{padding:'10px'}">
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
                                    <el-button style="border-radius:2px;" icon="Refresh" plain @click="getIngresss()">刷新</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="ingress-head-card" shadow="never" :body-style="{padding:'10px'}">
                        <el-row>
                            <el-col :span="2">
                                <div>
                                    <el-button style="border-radius:2px;" icon="Edit" type="primary" @click="createIngressDrawer = true" v-loading.fullscreen.lock="fullscreenLoading">创建</el-button>
                                </div>
                            </el-col>
                            <el-col :span="6">
                                <div>
                                    <el-input class="ingress-head-search" clearable placeholder="请输入" v-model="searchInput"></el-input>
                                    <el-button style="border-radius:2px;" icon="Search" type="primary" plain @click="getIngresss()">搜索</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="ingress-body-card" shadow="never" :body-style="{padding:'5px'}">
                        <el-table
                        style="width:100%;font-size:12px;margin-bottom:10px;"
                        :data="ingressList"
                        v-loading="appLoading">
                            <el-table-column width="10"></el-table-column>
                            <el-table-column align=left label="Ingress名">
                                <template v-slot="scope">
                                    <a class="ingress-body-ingressname">{{ scope.row.metadata.name }}</a>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="标签" min-width='120'>
                                <template v-slot="scope">
                                    <div v-for="(val, key) in scope.row.metadata.labels" :key="key">
                                        <el-popover
                                            placement="right"
                                            :width="200"
                                            trigger="hover"
                                            :content="key + ':' + val">
                                            <template #reference>
                                                <el-tag style="margin-bottom: 5px" type="warning">{{ ellipsis(key + ":" + val) }}</el-tag>
                                            </template>
                                        </el-popover>
                                    </div>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="Host" min-width='120'>
                                <template v-slot="scope">
                                    <div v-for="(item, index) in scope.row.spec.rules" :key="index">
                                        <el-popover
                                            placement="right"
                                            :width="200"
                                            trigger="hover"
                                            :content="item.host">
                                            <template #reference>
                                                <el-tag style="margin-bottom: 5px" type="danger">{{ ellipsis(item.host) }}</el-tag>
                                            </template>
                                        </el-popover>
                                    </div>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="Path">
                                <template v-slot="scope">
                                    <div v-for="(item, index) in scope.row.spec.rules" :key="index">
                                        <el-popover
                                            placement="right"
                                            :width="100"
                                            trigger="hover"
                                            :content="item.http.paths[0].path">
                                            <template #reference>
                                                <el-tag style="margin-bottom: 5px" type="danger">{{ item.http.paths[0].path }}</el-tag>
                                            </template>
                                        </el-popover>
                                    </div>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="EXTERNAL-IP">
                                <template v-slot="scope">
                                    <span>{{ scope.row.status.loadBalancer.ingress ? scope.row.status.loadBalancer.ingress[0].ip : '' }} </span>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="TLS">
                                <template v-slot="scope">
                                    <span>{{ scope.row.spec.tls ? 'YES' : '' }} </span>
                                </template>
                            </el-table-column>
                            <el-table-column align=center min-width="100" label="创建时间">
                                <template v-slot="scope">
                                    <el-tag type="info">{{ timeTrans(scope.row.metadata.creationTimestamp) }} </el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="操作" width="200">
                                <template v-slot="scope">
                                    <el-button size="small" style="border-radius:2px;" icon="Edit" type="primary" plain @click="getIngressDetail(scope)">YAML</el-button>
                                    <el-button size="small" style="border-radius:2px;" icon="Delete" type="danger" @click="handleConfirm(scope, '删除', delIngress)">删除</el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                        <el-pagination
                        class="ingress-body-pagination"
                        background
                        @size-change="handleSizeChange"
                        @current-change="handleCurrentChange"
                        :current-page="currentPage"
                        :page-sizes="pagesizeList"
                        :page-size="pagesize"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="ingressTotal">
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
                    <el-button type="primary" @click="updateIngress()">更 新</el-button>
                </span>
            </template>
        </el-dialog>
        <el-drawer
            v-model="createIngressDrawer"
            :direction="direction"
            :before-close="handleClose">
            <template #title>
                <h4>创建Ingress</h4>
            </template>
            <template #default>
                <el-row type="flex" justify="center">
                    <el-col :span="20">
                        <el-form ref="createIngress" :rules="createIngressRules" :model="createIngress" label-width="80px">
                            <el-form-item class="ingress-create-form" label="名称" prop="name">
                                <el-input v-model="createIngress.name"></el-input>
                            </el-form-item>
                            <el-form-item class="ingress-create-form" label="命名空间" prop="namespace">
                                <el-select v-model="createIngress.namespace" filterable placeholder="请选择">
                                    <el-option
                                    v-for="(item, index) in namespaceList"
                                    :key="index"
                                    :label="item.metadata.name"
                                    :value="item.metadata.name">
                                    </el-option>
                                </el-select>
                            </el-form-item>
                            <el-form-item class="SERVICE-create-form" label="标签" prop="label_str">
                                <el-input v-model="createIngress.label_str" placeholder="示例: project=ms,app=gateway"></el-input>
                            </el-form-item>
                            <el-form-item class="deploy-create-form" label="域名" prop="host">
                                <el-input v-model="createIngress.host" placeholder="示例: www.example.com"></el-input>
                            </el-form-item>
                            <el-form-item class="ingress-create-form" label="Path" prop="path">
                                <el-input v-model="createIngress.path" placeholder="示例: /abc"></el-input>
                            </el-form-item>
                            <el-form-item class="deploy-create-form" label="匹配类型" prop="path_type">
                                <el-select v-model="createIngress.path_type" placeholder="请选择">
                                    <el-option value="Prefix" label="Prefix"></el-option>
                                    <el-option value="Exact" label="Exact"></el-option>
                                    <el-option value="ImplementationSpecific" label="ImplementationSpecific"></el-option>
                                </el-select>
                            </el-form-item>
                            <el-form-item class="ingress-create-form" label="Service名" prop="service_name">
                                <el-input disabled v-model="createIngress.name"></el-input>
                            </el-form-item>
                            <el-form-item class="ingress-create-form" label="Service端口" prop="service_port">
                                <el-input v-model="createIngress.service_port" placeholder="示例: 80"></el-input>
                            </el-form-item>
                        </el-form>
                    </el-col>
                </el-row>
            </template>
            <template #footer>
                <el-button @click="createIngressDrawer = false">取消</el-button>
                <el-button type="primary" @click="submitForm('createIngress')">立即创建</el-button>
            </template>
        </el-drawer>
    </div>
</template>

<script>
import common from "../common/Config";
import httpClient from '../../utils/request';
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
            ingressList: [],
            ingressTotal: 0,
            getIngresssData: {
                url: common.k8sIngressList,
                params: {
                    filter_name: '',
                    namespace: '',
                    page: '',
                    limit: '',
                }
            },
            ingressDetail: {},
            getIngressDetailData: {
                url: common.k8sIngressDetail,
                params: {
                    ingress_name: '',
                    namespace: ''
                }
            },
            yamlDialog: false,
            updateIngressData: {
                url: common.k8sIngressUpdate,
                params: {
                    namespace: '',
                    content: ''
                }
            },
            delIngressData: {
                url: common.k8sIngressDel,
                params: {
                    ingress_name: '',
                    namespace: '',
                }
            },
            fullscreenLoading: false,
            direction: 'rtl',
            createIngressDrawer: false,
            createIngress: {
                name: '',
                namespace: '',
                label_str: '',
                host: '',
                path: '',
                path_type: '',
                service_name: '',
                service_port: '',
                hosts: {}
            },
            createIngressData: {
                url: common.k8sIngressCreate,
                params: {}
            },
            createIngressRules: {
                name: [{
                    required: true,
                    message: '请填写名称',
                    trigger: 'change'
                }],
                namespace: [{
                    required: true,
                    message: '请选择命名空间',
                    trigger: 'change'
                }],
                host: [{
                    required: true,
                    message: '请填写域名',
                    trigger: 'change'
                }],
                path: [{
                    required: true,
                    message: '请填写路径',
                    trigger: 'change'
                }],
                service_port: [{
                    required: true,
                    message: '请填写Service端口',
                    trigger: 'change'
                }],
                label_str: [{
                    required: true,
                    message: '请填写标签',
                    trigger: 'change'
                }],
                path_type: [{
                    required: true,
                    message: '请选择匹配类型',
                    trigger: 'change'
                }],
            }
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
            this.getIngresss()
        },
        handleCurrentChange(currentPage) {
            this.currentPage = currentPage;
            this.getIngresss()
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
        getIngresss() {
            this.appLoading = true
            this.getIngresssData.params.filter_name = this.searchInput
            this.getIngresssData.params.namespace = this.namespaceValue
            this.getIngresssData.params.page = this.currentPage
            this.getIngresssData.params.limit = this.pagesize
            httpClient.get(this.getIngresssData.url, {params: this.getIngresssData.params})
            .then(res => {
                this.ingressList = res.data.items
                this.ingressTotal = res.data.total
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.appLoading = false
        },
        getIngressDetail(e) {
            this.getIngressDetailData.params.ingress_name = e.row.metadata.name
            this.getIngressDetailData.params.namespace = this.namespaceValue
            httpClient.get(this.getIngressDetailData.url, {params: this.getIngressDetailData.params})
            .then(res => {
                this.ingressDetail = res.data
                this.contentYaml = this.transYaml(this.ingressDetail)
                this.yamlDialog = true
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        updateIngress() {
            let content = JSON.stringify(this.transObj(this.contentYaml))
            this.updateIngressData.params.namespace = this.namespaceValue
            this.updateIngressData.params.content = content
            httpClient.put(this.updateIngressData.url, this.updateIngressData.params)
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
        delIngress(e) {
            this.delIngressData.params.ingress_name = e.row.metadata.name
            this.delIngressData.params.namespace = this.namespaceValue
            httpClient.delete(this.delIngressData.url, {data: this.delIngressData.params})
            .then(res => {
                this.getIngresss()
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
        createIngressFunc() {
            let reg = new RegExp("(^[A-Za-z]+=[A-Za-z0-9]+).*")
            if (!reg.test(this.createIngress.label_str)) {
                this.$message.warning({
                    message: "标签填写异常，请确认后重新填写"
                })
                return
            }
            this.fullscreenLoading = true
            let label = new Map()
            let a = (this.createIngress.label_str).split(",")
            a.forEach(item => {
                let b = item.split("=")
                label[b[0]] = b[1]
            })
            let hosts = new Map()
            let httpPaths = []
            let httpPath = {
                path: this.createIngress.path,
                path_type: this.createIngress.path_type,
                service_name: this.createIngress.name,
                service_port: parseInt(this.createIngress.service_port)
            }
            httpPaths.push(httpPath)
            hosts[this.createIngress.host] = httpPaths

            this.createIngressData.params = this.createIngress
            this.createIngressData.params.label = label
            this.createIngressData.params.hosts = hosts
            httpClient.post(this.createIngressData.url, this.createIngressData.params)
            .then(res => {
                this.$message.success({
                message: res.msg
                })
                this.getIngresss()
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.resetForm('createIngress')
            this.fullscreenLoading = false
            this.createIngressDrawer = false
        },
        resetForm(formName) {
            this.$refs[formName].resetFields()
        },
        submitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.createIngressFunc()
                } else {
                    return false;
                }
            })
        }
    },
    watch: {
        namespaceValue: {
            handler() {
                localStorage.setItem('namespace', this.namespaceValue)
                this.currentPage = 1
                this.getIngresss()
            }
        },
    },
    beforeMount() {
        if (localStorage.getItem('namespace') !== undefined && localStorage.getItem('namespace') !== null) {
            this.namespaceValue = localStorage.getItem('namespace')
        }
        this.getNamespaces()
        this.getIngresss()
    }
}
</script>


<style scoped>
    .ingress-head-card,.ingress-body-card {
        border-radius: 1px;
        margin-bottom: 5px;
    }
    .ingress-head-search {
        width:160px;
        margin-right:10px; 
    }
    .ingress-body-ingressname {
        color: #4795EE;
    }
    .ingress-body-ingressname:hover {
        color: rgb(84, 138, 238);
        cursor: pointer;
        font-weight: bold;
    }
</style>