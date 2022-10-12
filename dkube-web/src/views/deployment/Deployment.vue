<template>
    <div class="deploy">
        <el-row>
            <el-col :span="24">
                <div>
                    <el-card class="deploy-head-card" shadow="never" :body-style="{padding:'10px'}">
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
                                    <el-button style="border-radius:2px;" icon="Refresh" plain @click="getDeployments()">刷新</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="deploy-head-card" shadow="never" :body-style="{padding:'10px'}">
                        <el-row>
                            <el-col :span="2">
                                <div>
                                    <el-button style="border-radius:2px;" icon="Edit" type="primary" @click="createDeploymentDrawer = true" v-loading.fullscreen.lock="fullscreenLoading">创建</el-button>
                                </div>
                            </el-col>
                            <el-col :span="6">
                                <div>
                                    <el-input class="deploy-head-search" clearable placeholder="请输入" v-model="searchInput"></el-input>
                                    <el-button style="border-radius:2px;" icon="Search" type="primary" plain @click="getDeployments()">搜索</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="deploy-body-card" shadow="never" :body-style="{padding:'5px'}">
                        <el-table
                        style="width:100%;font-size:12px;margin-bottom:10px;"
                        :data="deploymentList"
                        v-loading="appLoading">
                            <el-table-column width="20"></el-table-column>
                            <el-table-column align=left label="Deployment名">
                                <template v-slot="scope">
                                    <a class="deploy-body-deployname">{{ scope.row.metadata.name }}</a>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="标签">
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
                            <el-table-column align=center label="容器组">
                                <template v-slot="scope">
                                    <span>{{ scope.row.status.availableReplicas>0?scope.row.status.availableReplicas:0  }} / {{ scope.row.spec.replicas>0?scope.row.spec.replicas:0 }} </span>
                                </template>
                            </el-table-column>
                            <el-table-column align=center min-width="100" label="创建时间">
                                <template v-slot="scope">
                                    <el-tag type="info">{{ timeTrans(scope.row.metadata.creationTimestamp) }} </el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="镜像">
                                <template v-slot="scope">
                                    <div v-for="(val, key) in scope.row.spec.template.spec.containers" :key="key">
                                        <el-popover
                                            placement="right"
                                            :width="200"
                                            trigger="hover"
                                            :content="val.image">
                                            <template #reference>
                                                <el-tag style="margin-bottom: 5px">{{ ellipsis(val.image.split('/')[2]==undefined?val.image:val.image.split('/')[2]) }}</el-tag>
                                            </template>
                                        </el-popover>
                                    </div>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="操作" width="400">
                                <template v-slot="scope">
                                    <el-button size="small" style="border-radius:2px;" icon="Edit" type="primary" plain @click="getDeploymentDetail(scope)">YAML</el-button>
                                    <el-button size="small" style="border-radius:2px;" icon="Plus" type="primary" @click="handleScale(scope)">扩缩</el-button>
                                    <el-button size="small" style="border-radius:2px;" icon="RefreshLeft" type="primary" @click="handleConfirm(scope, '重启', restartDeployment)">重启</el-button>
                                    <el-button size="small" style="border-radius:2px;" icon="Delete" type="danger" @click="handleConfirm(scope, '删除', delDeployment)">删除</el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                        <el-pagination
                        class="deploy-body-pagination"
                        background
                        @size-change="handleSizeChange"
                        @current-change="handleCurrentChange"
                        :current-page="currentPage"
                        :page-sizes="pagesizeList"
                        :page-size="pagesize"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="deploymentTotal">
                        </el-pagination>
                    </el-card>
                </div>
            </el-col>
        </el-row>
        <el-drawer
            v-model="createDeploymentDrawer"
            :direction="direction"
            :before-close="handleClose">
            <template #title>
                <h4>创建Deployment</h4>
            </template>
            <template #default>
                <el-row type="flex" justify="center">
                    <el-col :span="20">
                        <el-form ref="createDeployment" :rules="createDeploymentRules" :model="createDeployment" label-width="80px">
                            <el-form-item class="deploy-create-form" label="名称" prop="name">
                                <el-input v-model="createDeployment.name"></el-input>
                            </el-form-item>
                            <el-form-item class="deploy-create-form" label="命名空间" prop="namespace">
                                <el-select v-model="createDeployment.namespace" filterable placeholder="请选择">
                                    <el-option
                                    v-for="(item, index) in namespaceList"
                                    :key="index"
                                    :label="item.metadata.name"
                                    :value="item.metadata.name">
                                    </el-option>
                                </el-select>
                            </el-form-item>
                            <el-form-item class="deploy-create-form" label="副本数" prop="replicas">
                                <el-input-number v-model="createDeployment.replicas" :min="1" :max="10"></el-input-number>
                                    <el-popover
                                        placement="top"
                                        :width="100"
                                        trigger="hover"
                                        content="申请副本数上限为10个">
                                        <template #reference>
                                            <el-icon style="width:2em;font-size:18px;color:#4795EE"><WarningFilled/></el-icon>
                                        </template>
                                    </el-popover>
                            </el-form-item>
                            <el-form-item class="deploy-create-form" label="镜像" prop="image">
                                <el-input v-model="createDeployment.image"></el-input>
                            </el-form-item>
                            <el-form-item class="deploy-create-form" label="标签" prop="label_str">
                                <el-input v-model="createDeployment.label_str" placeholder="示例: project=ms,app=gateway"></el-input>
                            </el-form-item>
                            <el-form-item class="deploy-create-form" label="资源配额" prop="resource">
                                <el-select v-model="createDeployment.resource" placeholder="请选择">
                                    <el-option value="0.5/1" label="0.5C1G"></el-option>
                                    <el-option value="1/2" label="1C2G"></el-option>
                                    <el-option value="2/4" label="2C4G"></el-option>
                                    <el-option value="4/8" label="4C8G"></el-option>
                                </el-select>
                            </el-form-item>
                            <el-form-item class="deploy-create-form" label="容器端口" prop="container_port">
                                <el-input v-model="createDeployment.container_port" placeholder="示例: 80"></el-input>
                            </el-form-item>
                            <el-form-item class="deploy-create-form" label="健康检查" prop="health">
                                <el-switch v-model="createDeployment.health_check" />
                            </el-form-item>
                            <el-form-item class="deploy-create-form" label="检查路径" prop="healthPath">
                                <el-input v-model="createDeployment.health_path" placeholder="示例: /health"></el-input>
                            </el-form-item>
                        </el-form>
                    </el-col>
                </el-row>
            </template>
            <template #footer>
                <el-button @click="createDeploymentDrawer = false">取消</el-button>
                <el-button type="primary" @click="submitForm('createDeployment')">立即创建</el-button>
            </template>
        </el-drawer>
        <el-dialog title="YAML信息" v-model="yamlDialog" width="45%" top="2%">
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
                    <el-button @click="this.yamlDialog = false">取 消</el-button>
                    <el-button type="primary" @click="updateDeployment()">更 新</el-button>
                </span>
            </template>
        </el-dialog>
        <el-dialog title="副本数调整" v-model="scaleDialog" width="25%">
            <div style="text-align:center">
                <span>实例数: </span>
                <el-input-number :step="1" v-model="scaleNum" :min="0" :max="30" label="描述文字"></el-input-number>
            </div>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="scaleDialog = false">取 消</el-button>
                    <el-button type="primary" @click="scaleDeployment()">更 新</el-button>
                </span>
            </template>
        </el-dialog>
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
            deploymentList: [],
            deploymentTotal: 0,
            getDeploymentsData: {
                url: common.k8sDeploymentList,
                params: {
                    filter_name: '',
                    namespace: '',
                    page: '',
                    limit: '',
                }
            },
            fullscreenLoading: false,
            direction: 'rtl',
            createDeploymentDrawer: false,
            createDeployment: {
                name: '',
                namespace: '',
                replicas: 1,
                image: '',
                resource: '',
                health_check: false,
                health_path: '',
                label_str: '',
                label: {},
                container_port: ''
            },
            createDeploymentData: {
                url: common.k8sDeploymentCreate,
                params: {}
            },
            createDeploymentRules: {
                name: [{
                    required: true,
                    message: '请填写名称',
                    trigger: 'change'
                }],
                image: [{
                    required: true,
                    message: '请填写镜像',
                    trigger: 'change'
                }],
                namespace: [{
                    required: true,
                    message: '请选择命名空间',
                    trigger: 'change'
                }],
                resource: [{
                    required: true,
                    message: '请选择配额',
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
            deploymentDetail: {},
            getDeploymentDetailData: {
                url: common.k8sDeploymentDetail,
                params: {
                    deployment_name: '',
                    namespace: ''
                }
            },
            yamlDialog: false,
            updateDeploymentData: {
                url: common.k8sDeploymentUpdate,
                params: {
                    namespace: '',
                    content: ''
                }
            },
            scaleNum: 0,
            scaleDialog: false,
            scaleDeploymentData: {
                url: common.k8sDeploymentScale,
                params: {
                    deployment_name: '',
                    namespace: '',
                    scale_num: ''
                }
            },
            restartDeploymentData: {
                url: common.k8sDeploymentRestart,
                params: {
                    deployment_name: '',
                    namespace: '',
                }
            },
            delDeploymentData: {
                url: common.k8sDeploymentDel,
                params: {
                    deployment_name: '',
                    namespace: '',
                }
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
            this.getDeployments()
        },
        handleCurrentChange(currentPage) {
            this.currentPage = currentPage;
            this.getDeployments()
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
        getDeployments() {
            this.appLoading = true
            this.getDeploymentsData.params.filter_name = this.searchInput
            this.getDeploymentsData.params.namespace = this.namespaceValue
            this.getDeploymentsData.params.page = this.currentPage
            this.getDeploymentsData.params.limit = this.pagesize
            httpClient.get(this.getDeploymentsData.url, {params: this.getDeploymentsData.params})
            .then(res => {
                this.deploymentList = res.data.items
                this.deploymentTotal = res.data.total
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.appLoading = false
        },
        getDeploymentDetail(e) {
            this.getDeploymentDetailData.params.deployment_name = e.row.metadata.name
            this.getDeploymentDetailData.params.namespace = this.namespaceValue
            httpClient.get(this.getDeploymentDetailData.url, {params: this.getDeploymentDetailData.params})
            .then(res => {
                this.deploymentDetail = res.data
                this.contentYaml = this.transYaml(this.deploymentDetail)
                this.yamlDialog = true
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        updateDeployment() {
            let content = JSON.stringify(this.transObj(this.contentYaml))
            this.updateDeploymentData.params.namespace = this.namespaceValue
            this.updateDeploymentData.params.content = content
            httpClient.put(this.updateDeploymentData.url, this.updateDeploymentData.params)
            .then(res => {
                this.$message.success({
                message: res.msg
                })
                this.getDeployments()
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.yamlDialog = false
        },

        handleScale(e) {
            this.scaleDialog = true
            this.deploymentDetail = e.row
            this.scaleNum = e.row.spec.replicas
        },

        scaleDeployment() {
            this.scaleDeploymentData.params.deployment_name = this.deploymentDetail.metadata.name
            this.scaleDeploymentData.params.namespace = this.namespaceValue
            this.scaleDeploymentData.params.scale_num = this.scaleNum
            httpClient.put(this.scaleDeploymentData.url, this.scaleDeploymentData.params)
            .then(res => {
                this.$message.success({
                message: res.msg
                })
                this.getDeployments()
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.scaleDialog = false
        },

        restartDeployment(e) {
            this.restartDeploymentData.params.deployment_name = e.row.metadata.name
            this.restartDeploymentData.params.namespace = this.namespaceValue
            httpClient.put(this.restartDeploymentData.url, this.restartDeploymentData.params)
            .then(res => {
                this.$message.success({
                message: res.msg
                })
                this.getDeployments()
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },

        delDeployment(e) {
            this.delDeploymentData.params.deployment_name = e.row.metadata.name
            this.delDeploymentData.params.namespace = this.namespaceValue
            httpClient.delete(this.delDeploymentData.url, {data: this.delDeploymentData.params})
            .then(res => {
                this.$message.success({
                message: res.msg
                })
                this.getDeployments()
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

        createDeployFunc() {
            let reg = new RegExp("(^[A-Za-z]+=[A-Za-z0-9]+).*")
            if (!reg.test(this.createDeployment.label_str)) {
                this.$message.warning({
                    message: "标签填写异常，请确认后重新填写"
                })
                return
            }

            this.fullscreenLoading = true
            let label = new Map()
            let cpu, memory
            let a = (this.createDeployment.label_str).split(",")
            a.forEach(item => {
                let b = item.split("=")
                label[b[0]] = b[1]
            })
            let resourceList = this.createDeployment.resource.split("/")
            cpu = resourceList[0]
            memory = resourceList[1] + "Gi"
            this.createDeploymentData.params = this.createDeployment
            this.createDeploymentData.params.container_port = parseInt(this.createDeployment.container_port)
            this.createDeploymentData.params.label = label
            this.createDeploymentData.params.cpu = cpu
            this.createDeploymentData.params.memory = memory
            httpClient.post(this.createDeploymentData.url, this.createDeploymentData.params)
            .then(res => {
                this.$message.success({
                message: res.msg
                })
                this.getDeployments()
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.resetForm('createDeployment')
            this.fullscreenLoading = false
            this.createDeploymentDrawer = false
        },
        resetForm(formName) {
            this.$refs[formName].resetFields()
        },
        submitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.createDeployFunc()
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
                this.getDeployments()
            }
        },
    },
    beforeMount() {
        if (localStorage.getItem('namespace') !== undefined && localStorage.getItem('namespace') !== null) {
            this.namespaceValue = localStorage.getItem('namespace')
        }
        this.getNamespaces()
        this.getDeployments()
    }
}
</script>


<style scoped>
    .deploy-head-card,.deploy-body-card {
        border-radius: 1px;
        margin-bottom: 5px;
    }
    .deploy-head-search {
        width:160px;
        margin-right:10px; 
    }
    .deploy-body-deployname {
        color: #4795EE;
    }
    .deploy-body-deployname:hover {
        color: rgb(84, 138, 238);
        cursor: pointer;
        font-weight: bold;
    }
</style>