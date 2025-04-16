<script lang='ts' setup>
import { ref } from 'vue'
import { useCarLogHook } from './index.tsx';
import { FormInstance } from 'element-plus';
import { formantHikDateTime, formatDateTime } from '@/utils/index.ts';
// import iconv from "iconv-lite";

defineOptions({
  name: "CarLog"
})

const formRef = ref<FormInstance>()
const {
  searchFormData, handleSearch, loading,
  handleCurrentChange, handleSizeChange,
  paginationData, handleRefresh,
  selectAreaValue, handleSelectArea,
  selectAreaOption,
  // dialogVisible, title, formData,
  // handleOpen,
  // handleClose, closeDialog,
  // operateAction, mpFormRules,
  tableData,
  // handleRow, selectUserOption
  // formRules, formRef, closeDialog, operateAction,handleClose
} = useCarLogHook();

const utf8ToGBK = (str: string) => {
  var encoder = new TextEncoder();
  var decoder = new TextDecoder('gbk');
  var uint8Array = encoder.encode(str);
  var str = decoder.decode(uint8Array);
  return str;
}
</script>

<template>
  <div class="app-container">
    <!-- 搜索框 -->
    <el-card shadow="never" class="search-wrapper">
      <el-form ref="searchFormRef" :inline="true" :model="searchFormData">
        <el-form-item prop="name" label="名称">
          <el-input v-model="searchFormData.carNum" placeholder="车牌号" />
        </el-form-item>
        <el-form-item prop="name" label="区域选择">
          <el-select clearable v-model="selectAreaValue" placeholder="选择区域筛选设备" style="width: 240px"
            @change="handleSelectArea">
            <el-option v-for="item in selectAreaOption" placeholder="请选择区域筛选" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleSearch">查询</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格 -->
    <el-card v-loading="loading" shadow="never">
      <!-- <div class="toolbar-wrapper">
        <el-button type="primary" plain icon="plus" @click="handleOpen('Add')">新增</el-button>
        <div>
          <el-tooltip content="刷新" effect="light">
            <el-button type="primary" icon="RefreshRight" circle plain @click="handleRefresh" />
          </el-tooltip>
        </div>
      </div> -->
      <div class="table-wrapper">
        <el-table :data="tableData" border stripe>
          <el-table-column prop="carNum" label="区域名称" min-width="200" />
          <el-table-column prop="uri" label="图片" min-width="200">
            <template #default="scope">
              <el-image style="width: 100px; height: 100px" :zoom-rate="1.2" :max-scale="7" :min-scale="0.2"
                :preview-src-list="['http://10.9.0.30:8888/preview/' + scope.row.uri]" show-progress :initial-index="4"
                :preview-teleported="true" fit="cover" :src="'http://10.9.0.30:8888/preview/' + scope.row.uri" />
            </template>
          </el-table-column>
          <el-table-column prop="device.area.name" label="抓拍区域" min-width="200"></el-table-column>
          <el-table-column prop="subTime" label="收到时间" min-width="200">
            <template #default="scope">
              {{ formantHikDateTime(scope.row.subTime) }}
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" min-width="200" />
          <el-table-column prop="createdAt" label="创建时间">
            <template #default="scope">
              {{ scope.row.createdAt }}
            </template>
          </el-table-column>
          <!-- <el-table-column label="操作" width="200">
            <template #default="scope">
              <el-button type="text" @click="handleRow(scope.row, 'Edit')">编辑</el-button>
              <el-button type="danger" @click="handleRow(scope.row, '')" text>删除</el-button>
            </template>
          </el-table-column> -->
        </el-table>
      </div>
      <div class="pager-wrapper">
        <el-pagination background :layout="paginationData.layout" :page-sizes="paginationData.pageSizes"
          :total="paginationData.total" :page-size="paginationData.pageSize" :currentPage="paginationData.currentPage"
          @size-change="handleSizeChange" @current-change="handleCurrentChange" />
      </div>
    </el-card>


  </div>
</template>

<style scoped></style>
