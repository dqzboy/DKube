<template>
    <div class="main-body-div">
        <el-row>
            <el-col :span="24">
                <div>
                    <img class="main-body-img" src="../../assets/img/403.png" />
                </div>
                </el-col>
            <el-col :span="24">
                <div>
                    <p class="status-code">403</p>
                    <p class="status-describe">你暂时无权限访问该页面······</p>
                </div>
            </el-col>
        </el-row>
    </div>
</template>

<script>
export default {
}
</script>


<style scoped>
  .main-body-img {
    margin-top: 150px
  }
  .main-body-div {
    text-align: center;
    height: 100vh;
    width: 100vw;
  }
  .status-code {
    margin-top: 20px;
    margin-bottom: 10px;
    font-size: 95px;
    font-weight: bold;
    color: rgb(54, 95, 230);
  }
  .status-describe {
    color: rgb(145, 143, 143);
  }
</style>