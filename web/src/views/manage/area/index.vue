<script lang='ts' setup>
import { ref } from 'vue';
import { useAreaHook } from './index.tsx';
import { FormInstance } from 'element-plus';
import { formatDateTime } from '@/utils/index.ts';

defineOptions({
  name: "Area"
})

const formRef = ref<FormInstance>()
const {
  searchFormData, handleSearch, loading,
  handleCurrentChange, handleSizeChange,
  paginationData, handleRefresh,
  dialogVisible, title, formData,
  handleOpen,
  handleClose, closeDialog,
  operateAction, mpFormRules,
  tableData,handleRow
  // formRules, formRef, closeDialog, operateAction,handleClose
} = useAreaHook();
</script>

<template>
  <div class="app-container">
    <!-- 搜索框 -->
    <el-card shadow="never" class="search-wrapper">
      <el-form ref="searchFormRef" :inline="true" :model="searchFormData">
        <el-form-item prop="name" label="名称">
          <el-input v-model="searchFormData.name" placeholder="名称" />
        </el-form-item>
        <el-form-item prop="remark" label="备注">
          <el-input v-model="searchFormData.remark" placeholder="备注" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleSearch">查询</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格 -->
    <el-card v-loading="loading" shadow="never">
      <div class="toolbar-wrapper">
        <el-button type="primary" plain icon="plus" @click="handleOpen('Add')">新增</el-button>
        <div>
          <el-tooltip content="刷新" effect="light">
            <el-button type="primary" icon="RefreshRight" circle plain @click="handleRefresh" />
          </el-tooltip>
        </div>
      </div>
      <div class="table-wrapper">
        <el-table :data="tableData" border stripe>
          <el-table-column prop="name" label="区域名称" min-width="200" />
          <el-table-column prop="remark" label="备注" min-width="200" />
          <el-table-column prop="createdAt" label="创建时间">
            <template #default="scope">
                {{ formatDateTime(scope.row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200">
            <template #default="scope">
              <el-button type="text"  @click="handleRow(scope.row,'Edit')" >编辑</el-button>
              <el-button type="danger" @click="handleRow(scope.row,'')"  text>删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div class="pager-wrapper">
        <el-pagination background :layout="paginationData.layout" :page-sizes="paginationData.pageSizes"
          :total="paginationData.total" :page-size="paginationData.pageSize" :currentPage="paginationData.currentPage"
          @size-change="handleSizeChange" @current-change="handleCurrentChange" />
      </div>
    </el-card>

    <!-- 编辑添加 -->
    <el-dialog v-model="dialogVisible" :title="title" @closed="handleClose(formRef!)" width="30%">
      <el-form ref="formRef" :model="formData" :rules="mpFormRules" label-width="100px" label-position="left"
        style="width: 95%; margin-top: 15px">
        <el-form-item label="区域名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入区域名称" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="formData.remark" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取消</el-button>
          <el-button type="primary" @click="operateAction(formRef!)">确认</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped></style>
