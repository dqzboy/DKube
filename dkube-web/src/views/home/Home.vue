<template>
  <div class="home">
      <el-collapse v-model="activeNames">
          <el-collapse-item title="集群资源" name="1">
              <el-row :gutter="10" style="margin-bottom: 10px;">
                  <el-col :span="5">
                      <el-card class="home-node-card" :body-style="{padding:'10px'}">
                          <div style="float:left;padding-top:20%">
                              <el-progress  :stroke-width="20" :show-text="false" type="circle" :percentage="namespaceActive/namespaceTotal * 100"></el-progress>
                          </div>
                          <div>
                              <p class="home-node-card-title">命名空间: Active/总量</p>
                              <p class="home-node-card-num">{{ namespaceActive }}/{{ namespaceTotal }}</p>
                          </div>
                      </el-card>
                  </el-col>
                  <el-col :span="5">
                      <el-card class="home-node-card" :body-style="{padding:'10px'}">
                          <div>
                              <p class="home-node-card-title">服务数</p>
                              <p class="home-node-card-num">{{ deploymentTotal }}</p>
                          </div>
                      </el-card>
                  </el-col>
                  <el-col :span="5">
                      <el-card class="home-node-card" :body-style="{padding:'10px'}">
                          <div>
                              <p class="home-node-card-title">实例数</p>
                              <p class="home-node-card-num">{{ podTotal }}</p>
                          </div>
                      </el-card>
                  </el-col>
              </el-row>
          </el-collapse-item>
          <el-collapse-item title="节点资源" name="2">
              <el-row :gutter="10" style="margin-bottom: 10px;">
                  <el-col :span="5">
                      <el-card class="home-node-card" :body-style="{padding:'10px'}">
                          <div style="float:left;padding-top:20%">
                              <el-progress :stroke-width="20" :show-text="false" type="circle" :percentage="nodeTotal/nodeTotal * 100"></el-progress>
                          </div>
                          <div>
                              <p class="home-node-card-title">节点: Ready/总数量</p>
                              <p class="home-node-card-num">{{ nodeTotal }}/{{ nodeTotal }}</p>
                          </div>
                      </el-card>
                  </el-col>
                  <el-col :span="5">
                      <el-card class="home-node-card" :body-style="{padding:'10px'}">
                          <div style="float:left;padding-top:20%">
                              <el-progress :stroke-width="20" :show-text="false" type="circle" :percentage="nodeCpuAllocatable/nodeCpuCapacity * 100"></el-progress>
                          </div>
                          <div>
                              <p class="home-node-card-title">CPU: 可分配/容量</p>
                              <p class="home-node-card-num">{{ nodeCpuAllocatable }}/{{ nodeCpuCapacity }}</p>
                          </div>
                      </el-card>
                  </el-col>
                  <el-col :span="5">
                      <el-card class="home-node-card" :body-style="{padding:'10px'}">
                          <div style="float:left;padding-top:20%">
                              <el-progress :stroke-width="20" :show-text="false" type="circle" :percentage="nodeMemAllocatable/nodeMemCapacity * 100"></el-progress>
                          </div>
                          <div>
                              <p class="home-node-card-title">内存: 可分配/容量</p>
                              <p class="home-node-card-num">{{ specTrans(nodeMemAllocatable) }}Gi/{{ specTrans(nodeMemCapacity) }}Gi</p>
                          </div>
                      </el-card>
                  </el-col>
                  <el-col :span="5">
                      <el-card class="home-node-card" :body-style="{padding:'10px'}">
                          <div style="float:left;padding-top:20%">
                              <el-progress :stroke-width="20" :show-text="false" type="circle" :percentage="nodePodAllocatable/nodePodAllocatable * 100"></el-progress>
                          </div>
                          <div>
                              <p class="home-node-card-title">POD: 可分配/容量</p>
                              <p class="home-node-card-num">{{ nodePodAllocatable }}/{{ nodePodAllocatable }}</p>
                          </div>
                      </el-card>
                  </el-col>
              </el-row>
          </el-collapse-item>
          <el-collapse-item title="资源统计" name="3">
              <el-row :gutter="10">
                  <el-col :span="24" style="margin-bottom: 10px;">
                      <el-card class="home-dash-card" :body-style="{padding:'10px'}">
                          <div id="podNumDash" style="height: 300px;">
                          </div>
                      </el-card>
                  </el-col>
                  <el-col :span="24">
                      <el-card class="home-dash-card" :body-style="{padding:'10px'}">
                          <div id="deployNumDash" style="height: 300px;">
                          </div>
                      </el-card>
                  </el-col>
              </el-row>
          </el-collapse-item>
      </el-collapse>
  </div>
</template>

<script>
import * as echarts from 'echarts'
import common from "../common/Config";
import httpClient from '../../utils/request';
export default {
  data() {
      return {
          activeNames: ["1", "2", "3"],
          namespaceActive: 0,
          namespaceValue: 'default',
          namespaceTotal: 0,
          namespaceListUrl: common.k8sNamespaceList,
          nodeTotal: 0,
          nodeCpuAllocatable: 0,
          nodeCpuCapacity: 0,
          nodeMemAllocatable: 0,
          nodeMemCapacity: 0,
          nodePodAllocatable: 0,
          nodePodCapacity: 0,
          getNodesData: {
              url: common.k8sNodeList,
              params: {}
          },
          deploymentTotal: 0,
          getDeploymentsData: {
              url: common.k8sDeploymentList,
              params: {
                  namespace: '',
              }
          },
          podTotal: 0,
          getPodsData: {
              url: common.k8sPodList,
              params: {
                  namespace: '',
              }
          },
          podNumNp: [],
          podNumNpUrl: common.k8sPodNumNp,
          deploymentNumNp: [],
          deploymentNumNpUrl: common.k8sDeploymentNumNp
      }
  },
  methods: {
      getNamespaces() {
          httpClient.get(this.namespaceListUrl)
          .then(res => {
              this.namespaceTotal = res.data.total
              let namespaceList = res.data.items
              let index
              for (index in namespaceList) {
                  if (namespaceList[index].status.phase === "Active" ) {
                      this.namespaceActive = this.namespaceActive + 1
                  }
              }
          })
          .catch(res => {
              this.$message.error({
              message: res.msg
              })
          })
      },
      specTrans(num) {
          let a = num / 1024 / 1024
          return a.toFixed(0)
      },
      getNodes() {
          httpClient.get(this.getNodesData.url, {params: this.getNodesData.params})
          .then(res => {
              this.nodeTotal = res.data.total
              let nodeList = res.data.items
              let index
              for (index in nodeList) {
                  let isnum = /^\d+$/.test(nodeList[index].status.allocatable.cpu);
                  if (!isnum) {
                      continue
                  }
                  this.nodeCpuAllocatable = parseInt(nodeList[index].status.allocatable.cpu) + this.nodeCpuAllocatable
                  this.nodeCpuCapacity = parseInt(nodeList[index].status.capacity.cpu) + this.nodeCpuCapacity
                  this.nodeMemAllocatable = parseInt(nodeList[index].status.allocatable.memory) + this.nodeMemAllocatable
                  this.nodeMemCapacity = parseInt(nodeList[index].status.capacity.memory) + this.nodeMemCapacity
                  this.nodePodAllocatable = parseInt(nodeList[index].status.allocatable.pods) + this.nodePodAllocatable
                  this.nodePodCapacity = parseInt(nodeList[index].status.capacity.pods) + this.nodePodCapacity
              }
          })
          .catch(res => {
              this.$message.error({
              message: res.msg
              })
          })
      },
      getDeployments() {
          this.getDeploymentsData.params.namespace = this.namespaceValue
          httpClient.get(this.getDeploymentsData.url, {params: this.getDeploymentsData.params})
          .then(res => {
              this.deploymentTotal = res.data.total
          })
          .catch(res => {
              this.$message.error({
              message: res.msg
              })
          })
      },
      getPods() {
          this.getPodsData.params.namespace = this.namespaceValue
          httpClient.get(this.getPodsData.url, {params: this.getPodsData.params})
          .then(res => {
              this.podTotal = res.data.total
          })
          .catch(res => {
              this.$message.error({
              message: res.msg
              })
          })
      },
      getDeploymentNumNp() {
          httpClient.get(this.deploymentNumNpUrl)
          .then(res => {
              this.deploymentNumNp = res.data
              this.getDeployNumDash()
          })
          .catch(res => {
              this.$message.error({
              message: res.msg
              })
          })
      },
      getPodNumNp() {
          httpClient.get(this.podNumNpUrl)
          .then(res => {
              this.podNumNp = res.data
              this.getPodNumDash()
          })
          .catch(res => {
              this.$message.error({
              message: res.msg
              })
          })
      },
      getPodNumDash(){
          if (this.podNumDash != null && this.podNumDash != "" && this.podNumDash != undefined) {
              this.podNumDash.dispose()
          }
          this.podNumDash = echarts.init(document.getElementById('podNumDash'));
          this.podNumDash.setOption({
              title: { text: 'Pods per Namespace', textStyle: {color:'rgb(134, 135, 136)'}},
              color: ['#67E0E3', '#9FE6B8', '#FFDB5C','#ff9f7f', '#fb7293', '#E062AE', '#E690D1', '#e7bcf3', '#9d96f5', '#8378EA', '#96BFFF'],
              tooltip: { 
                  trigger: "axis", 
                  axisPointer: { 
                      type: "cross", 
                      label: { 
                          backgroundColor: "#76baf1" 
                      } 
                  } 
              },
              legend: {
                  data: ['Pods']
              },
              dataset: {
                  dimensions: ['namespace','pod_num'],
                  source: this.podNumNp
              },
              xAxis: {
                  type: 'category',
                  axisLabel:{
                      interval: 0,
                      formatter: function (value) {
                          return value.length>5?value.substring(0,5)+'...':value
                      }
                  },
              },
              yAxis: [
                  {type: 'value'}
              ],
              series: [{
                  name: 'Pods',
                  type: 'bar',
                  label: {
                      show: true,
                      position: 'top'
                      }
                  }
              ]
          });
      },
      getDeployNumDash(){
          if (this.deployNumDash != null && this.deployNumDash != "" && this.deployNumDash != undefined) {
              this.deployNumDash.dispose()
          }
          this.deployNumDash = echarts.init(document.getElementById('deployNumDash'));

          this.deployNumDash.setOption({
              title: { text: 'Deployments per Namespace', textStyle: {color:'rgb(134, 135, 136)'}},
              color: ['#9FE6B8', '#FFDB5C','#ff9f7f', '#fb7293', '#E062AE', '#E690D1', '#e7bcf3', '#9d96f5', '#8378EA', '#96BFFF'],
              tooltip: { trigger: "axis", axisPointer: { type: "cross", label: { backgroundColor: "#76baf1" } } },
              legend: {
                  data: ['Deployments']
              },
              dataset: {
                  dimensions: ['namespace','deployment_num'],
                  source: this.deploymentNumNp
              },
              xAxis: {
                  type: 'category',
                  axisLabel:{
                      interval: 0,
                      formatter: function (value) {
                          return value.length>5?value.substring(0,5)+'...':value
                      }
                  },
              },
              yAxis: [
                  {type: 'value'}
              ],
              series: [{
                  name: 'Deployments',
                  type: 'bar',
                  label: {
                      show: true,
                      position: 'top'
                      }
                  }
              ]
          });
      },
  },
  beforeMount() {
      this.getNamespaces()
      this.getNodes()
      this.getDeployments()
      this.getPods()
      this.getDeploymentNumNp()
      this.getPodNumNp()
  }
}
</script>

<style scoped>
  :v-deep .el-collapse-item__header {
      font-size: 16px;
  }
  .home-node-card {
      border-radius:1px;
      text-align: center;
      background-color: rgb(250, 253, 255);
  }
  .home-dash-card {
      border-radius:1px;
  }
  .home-node-card-title {
      font-size: 12px;
  }
  .home-node-card-num {
      font-size: 22px;
      font-weight: bold;
      color: rgb(63, 92, 135);
  }
  :v-deep .el-progress-circle {
      height: 50px !important;
      width: 50px !important;
  }
</style>