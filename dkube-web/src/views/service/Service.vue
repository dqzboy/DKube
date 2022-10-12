<template>
    <div class="service">
        <el-row>
            <el-col :span="24">
                <div>
                    <el-card class="service-head-card" shadow="never" :body-style="{padding:'10px'}">
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
                                    <el-button style="border-radius:2px;" icon="Refresh" plain @click="getServices()">刷新</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="service-head-card" shadow="never" :body-style="{padding:'10px'}">
                        <el-row>
                            <el-col :span="2">
                                <div>
                                    <el-button style="border-radius:2px;" icon="Edit" type="primary" @click="createServiceDrawer = true" v-loading.fullscreen.lock="fullscreenLoading">创建</el-button>
                                </div>
                            </el-col>
                            <el-col :span="6">
                                <div>
                                    <el-input class="service-head-search" clearable placeholder="请输入" v-model="searchInput"></el-input>
                                    <el-button style="border-radius:2px;" icon="Search" type="primary" plain @click="getServices()">搜索</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="service-body-card" shadow="never" :body-style="{padding:'5px'}">
                        <el-table
                        style="width:100%;font-size:12px;margin-bottom:10px;"
                        :data="serviceList"
                        v-loading="appLoading">
                            <el-table-column width="20"></el-table-column>
                            <el-table-column align=left label="Service名">
                                <template v-slot="scope">
                                    <a class="service-body-servicename">{{ scope.row.metadata.name }}</a>
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
                            <el-table-column align=center label="类型">
                                <template v-slot="scope">
                                    <span style="font-weight:bold;">{{ scope.row.spec.type }} </span>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="CLUSTER-IP">
                                <template v-slot="scope">
                                    <span>{{ scope.row.spec.clusterIP }} </span>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="EXTERNAL-IP">
                                <template v-slot="scope">
                                    <span>{{ scope.row.status.loadBalancer.ingress ? scope.row.status.loadBalancer.ingress[0].ip : '' }} </span>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="端口">
                                <template v-slot="scope">
                                    <span v-if="!scope.row.spec.ports[0].nodePort">{{ scope.row.spec.ports[0].port }}/{{ scope.row.spec.ports[0].protocol }}</span>
                                    <span v-if="scope.row.spec.ports[0].nodePort">{{ scope.row.spec.ports[0].port }}:{{ scope.row.spec.ports[0].nodePort }}/{{ scope.row.spec.ports[0].protocol }}</span>
                                </template>
                            </el-table-column>
                            <el-table-column align=center min-width="100" label="创建时间">
                                <template v-slot="scope">
                                    <el-tag type="info">{{ timeTrans(scope.row.metadata.creationTimestamp) }} </el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="操作" width="200">
                                <template v-slot="scope">
                                    <el-button size="small" style="border-radius:2px;" icon="Edit" type="primary" plain @click="getServiceDetail(scope)">YAML</el-button>
                                    <el-button size="small" style="border-radius:2px;" icon="Delete" type="danger" @click="handleConfirm(scope, '删除', delService)">删除</el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                        <el-pagination
                        class="service-body-pagination"
                        background
                        @size-change="handleSizeChange"
                        @current-change="handleCurrentChange"
                        :current-page="currentPage"
                        :page-sizes="pagesizeList"
                        :page-size="pagesize"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="serviceTotal">
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
                    <el-button type="primary" @click="updateService()">更 新</el-button>
                </span>
            </template>
        </el-dialog>
        <el-drawer
            v-model="createServiceDrawer"
            :direction="direction"
            :before-close="handleClose">
            <template #title>
                <h4>创建Service</h4>
            </template>
            <template #default>
                <el-row type="flex" justify="center">
                    <el-col :span="20">
                        <el-form ref="createService" :rules="createServiceRules" :model="createService" label-width="80px">
                            <el-form-item class="service-create-form" label="名称" prop="name">
                                <el-input v-model="createService.name"></el-input>
                            </el-form-item>
                            <el-form-item class="service-create-form" label="命名空间" prop="namespace">
                                <el-select v-model="createService.namespace" filterable placeholder="请选择">
                                    <el-option
                                    v-for="(item, index) in namespaceList"
                                    :key="index"
                                    :label="item.metadata.name"
                                    :value="item.metadata.name">
                                    </el-option>
                                </el-select>
                            </el-form-item>
                            <el-form-item class="service-create-form" label="类型" prop="type">
                                <el-select v-model="createService.type" placeholder="请选择">
                                    <el-option value="ClusterIP" label="ClusterIP"></el-option>
                                    <el-option value="NodePort" label="NodePort"></el-option>
                                </el-select>
                            </el-form-item>
                            <el-form-item class="deploy-create-form" label="容器端口" prop="container_port">
                                <el-input v-model="createService.container_port" placeholder="示例: 80"></el-input>
                            </el-form-item>
                            <el-form-item class="service-create-form" label="Service端口" prop="port">
                                <el-input v-model="createService.port" placeholder="示例: 80"></el-input>
                            </el-form-item>
                            <el-form-item v-if="createService.type == 'NodePort'" class="service-create-form" label="NodePort" prop="node_port">
                                <el-input v-model="createService.node_port" placeholder="示例: 30001"></el-input>
                            </el-form-item>
                            <el-form-item class="SERVICE-create-form" label="标签" prop="label_str">
                                <el-input v-model="createService.label_str" placeholder="示例: project=ms,app=gateway"></el-input>
                            </el-form-item>
                        </el-form>
                    </el-col>
                </el-row>
            </template>
            <template #footer>
                <el-button @click="createServiceDrawer = false">取消</el-button>
                <el-button type="primary" @click="submitForm('createService')">立即创建</el-button>
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
            serviceList: [],
            serviceTotal: 0,
            getServicesData: {
                url: common.k8sServiceList,
                params: {
                    filter_name: '',
                    namespace: '',
                    page: '',
                    limit: '',
                }
            },
            serviceDetail: {},
            getServiceDetailData: {
                url: common.k8sServiceDetail,
                params: {
                    service_name: '',
                    namespace: ''
                }
            },
            yamlDialog: false,
            updateServiceData: {
                url: common.k8sServiceUpdate,
                params: {
                    namespace: '',
                    content: ''
                }
            },
            delServiceData: {
                url: common.k8sServiceDel,
                params: {
                    service_name: '',
                    namespace: '',
                }
            },
            fullscreenLoading: false,
            direction: 'rtl',
            createServiceDrawer: false,
            createService: {
                name: '',
                namespace: '',
                type: 'ClusterIP',
                container_port: '',
                port: '',
                node_port: '',
                label: {},
                label_str: ''
            },
            createServiceData: {
                url: common.k8sServiceCreate,
                params: {}
            },
            createServiceRules: {
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
                port: [{
                    required: true,
                    message: '请填写Service端口',
                    trigger: 'change'
                }],
                node_port: [{
                    required: true,
                    message: '请填写NodePort',
                    trigger: 'change'
                }],
                label_str: [{
                    required: true,
                    message: '请填写标签',
                    trigger: 'change'
                }],
                container_port: [{
                    required: true,
                    message: '请填写容器端口',
                    trigger: 'change'
                }],
            },
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
            this.getServices()
        },
        handleCurrentChange(currentPage) {
            this.currentPage = currentPage;
            this.getServices()
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
        getServices() {
            this.appLoading = true
            this.getServicesData.params.filter_name = this.searchInput
            this.getServicesData.params.namespace = this.namespaceValue
            this.getServicesData.params.page = this.currentPage
            this.getServicesData.params.limit = this.pagesize
            httpClient.get(this.getServicesData.url, {params: this.getServicesData.params})
            .then(res => {
                this.serviceList = res.data.items
                this.serviceTotal = res.data.total
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.appLoading = false
        },
        getServiceDetail(e) {
            this.getServiceDetailData.params.service_name = e.row.metadata.name
            this.getServiceDetailData.params.namespace = this.namespaceValue
            httpClient.get(this.getServiceDetailData.url, {params: this.getServiceDetailData.params})
            .then(res => {
                this.serviceDetail = res.data
                this.contentYaml = this.transYaml(this.serviceDetail)
                this.yamlDialog = true
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        updateService() {
            let content = JSON.stringify(this.transObj(this.contentYaml))
            this.updateServiceData.params.namespace = this.namespaceValue
            this.updateServiceData.params.content = content
            httpClient.put(this.updateServiceData.url, this.updateServiceData.params)
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
        delService(e) {
            this.delServiceData.params.service_name = e.row.metadata.name
            this.delServiceData.params.namespace = this.namespaceValue
            httpClient.delete(this.delServiceData.url, {data: this.delServiceData.params})
            .then(res => {
                this.getServices()
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
        createServiceFunc() {
            let reg = new RegExp("(^[A-Za-z]+=[A-Za-z0-9]+).*")
            if (!reg.test(this.createService.label_str)) {
                this.$message.warning({
                    message: "标签填写异常，请确认后重新填写"
                })
                return
            }
            this.fullscreenLoading = true
            let label = new Map()
            let a = (this.createService.label_str).split(",")
            a.forEach(item => {
                let b = item.split("=")
                label[b[0]] = b[1]
            })
            this.createServiceData.params = this.createService
            this.createServiceData.params.label = label
            this.createServiceData.params.container_port = parseInt(this.createService.container_port)
            this.createServiceData.params.port = parseInt(this.createService.port)
            this.createServiceData.params.node_port = parseInt(this.createService.node_port)
            httpClient.post(this.createServiceData.url, this.createServiceData.params)
            .then(res => {
                this.$message.success({
                message: res.msg
                })
                this.getServices()
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.resetForm('createService')
            this.fullscreenLoading = false
            this.createServiceDrawer = false
        },
        resetForm(formName) {
            this.$refs[formName].resetFields()
        },
        submitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.createServiceFunc()
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
                this.getServices()
            }
        },
    },
    beforeMount() {
        if (localStorage.getItem('namespace') !== undefined && localStorage.getItem('namespace') !== null) {
            this.namespaceValue = localStorage.getItem('namespace')
        }
        this.getNamespaces()
        this.getServices()
    }
}
</script>


<style scoped>
    .service-head-card,.service-body-card {
        border-radius: 1px;
        margin-bottom: 5px;
    }
    .service-head-search {
        width:160px;
        margin-right:10px; 
    }
    .service-body-servicename {
        color: #4795EE;
    }
    .service-body-servicename:hover {
        color: rgb(84, 138, 238);
        cursor: pointer;
        font-weight: bold;
    }
</style>