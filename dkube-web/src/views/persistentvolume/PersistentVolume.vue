<template>
    <div class="pv">
        <el-row>
            <el-col :span="24">
                <div>
                    <el-card class="pv-head-card" shadow="never" :body-style="{padding:'10px'}">
                        <el-row>
                            <el-col :span="2">
                                <div>
                                    <el-button disabled style="border-radius:2px;" icon="Edit" type="primary">创建</el-button>
                                </div>
                            </el-col>
                            <el-col :span="6">
                                <div>
                                    <el-input class="pv-head-search" clearable placeholder="请输入" v-model="searchInput"></el-input>
                                    <el-button style="border-radius:2px;" icon="Search" type="primary" plain @click="getPvs()">搜索</el-button>
                                </div>
                            </el-col>
                            <el-col :span="2" :offset="14">
                                <div>
                                    <el-button style="border-radius:2px;" icon="Refresh" plain @click="getPvs()">刷新</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="pv-body-card" shadow="never" :body-style="{padding:'5px'}">
                        <el-table
                        style="width:100%;font-size:12px;margin-bottom:10px;"
                        :data="pvList"
                        v-loading="appLoading">
                            <el-table-column width="20"></el-table-column>
                            <el-table-column align=left label="PV名">
                                <template v-slot="scope">
                                    <a class="pv-body-pvname">{{ scope.row.metadata.name }}</a>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="状态">
                                <template v-slot="scope">
                                    <span :class="[scope.row.status.phase === 'Bound' ? 'success-status' : 'error-status']">{{ scope.row.status.phase }}</span>
                                </template>
                            </el-table-column>
                            <el-table-column align=center prop="spec.accessModes[0]" label="访问模式"></el-table-column>
                            <el-table-column align=center prop="spec.capacity.storage" label="容量"></el-table-column>
                            <el-table-column align=center prop="spec.claimRef.name" label="PVC"></el-table-column>
                            <el-table-column align=center min-width="100" label="创建时间">
                                <template v-slot="scope">
                                    <el-tag type="info">{{ timeTrans(scope.row.metadata.creationTimestamp) }} </el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="操作" min-width="120">
                                <template v-slot="scope">
                                    <el-button size="small" style="border-radius:2px;" icon="Edit" type="primary" plain @click="getPvDetail(scope)">YAML</el-button>
                                    <el-button size="small" style="border-radius:2px;" icon="Delete" type="danger" @click="handleConfirm(scope, '删除', delPv)">删除</el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                        <el-pagination
                        class="pv-body-pagination"
                        background
                        @size-change="handleSizeChange"
                        @current-change="handleCurrentChange"
                        :current-page="currentPage"
                        :page-sizes="pagesizeList"
                        :page-size="pagesize"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="pvTotal">
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
                    <el-button disabled type="primary" @click="updatePv()">更 新</el-button>
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
            pvList: [],
            pvTotal: 0,
            getPvsData: {
                url: common.k8sPvList,
                params: {
                    filter_name: '',
                    namespace: '',
                    page: '',
                    limit: '',
                }
            },
            pvDetail: {},
            getPvDetailData: {
                url: common.k8sPvDetail,
                params: {
                    pv_name: '',
                    namespace: ''
                }
            },
            yamlDialog: false,
            updatePvData: {
                url: common.k8sPvUpdate,
                params: {
                    namespace: '',
                    content: ''
                }
            },
            delPvData: {
                url: common.k8spvDel,
                params: {
                    pv_name: '',
                    namespace: '',
                }
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
            this.getPvs()
        },
        handleCurrentChange(currentPage) {
            this.currentPage = currentPage;
            this.getPvs()
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
        specTrans(str) {
            if ( str.indexOf('Ki') == -1 ) {
                return str
            }
            let num = str.slice(0,-2) / 1024 / 1024
            return num.toFixed(0)
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
        getPvs() {
            this.appLoading = true
            this.getPvsData.params.filter_name = this.searchInput
            this.getPvsData.params.namespace = this.namespaceValue
            this.getPvsData.params.page = this.currentPage
            this.getPvsData.params.limit = this.pagesize
            httpClient.get(this.getPvsData.url, {params: this.getPvsData.params})
            .then(res => {
                this.pvList = res.data.items
                this.pvTotal = res.data.total
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.appLoading = false
        },
        getPvDetail(e) {
            this.getPvDetailData.params.pv_name = e.row.metadata.name
            this.getPvDetailData.params.namespace = this.namespaceValue
            httpClient.get(this.getPvDetailData.url, {params: this.getPvDetailData.params})
            .then(res => {
                this.pvDetail = res.data
                this.contentYaml = this.transYaml(this.pvDetail)
                this.yamlDialog = true
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        updatePv() {
            let content = JSON.stringify(this.transObj(this.contentYaml))
            this.updatePvData.params.namespace = this.namespaceValue
            this.updatePvData.params.content = content
            httpClient.put(this.updatePvData.url, this.updatePvData.params)
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
        delPv(e) {
            this.delPvData.params.pv_name = e.row.metadata.name
            this.delPvData.params.namespace = this.namespaceValue
            httpClient.delete(this.delPvData.url, {data: this.delPvData.params})
            .then(res => {
                this.getPvs()
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
    },
    beforeMount() {
        this.getPvs()
    }
}
</script>


<style scoped>
    .pv-head-card,.pv-body-card {
        border-radius: 1px;
        margin-bottom: 5px;
    }
    .pv-head-search {
        width:160px;
        margin-right:10px; 
    }
    .pv-body-pvname {
        color: #4795EE;
    }
    .pv-body-pvname:hover {
        color: rgb(84, 138, 238);
        cursor: pointer;
        font-weight: bold;
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
</style>