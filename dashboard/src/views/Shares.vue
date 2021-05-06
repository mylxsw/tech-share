<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2">
                <b-card-text>
                    <b-form inline @submit="searchSubmit">
                        <b-form-select v-model="search.status" class="mb-2 mr-sm-2 mb-sm-0" placeholder="状态" :options="status_options"></b-form-select>
                        <b-button variant="light" type="submit">搜索</b-button>
                    </b-form>   
                </b-card-text>
            </b-card>

            <b-card class="mb-2" border-variant="warning" header="分享列表" header-bg-variant="warning">
                <b-card-body>
                    <b-table :items="shares" :fields="shares_fields">
                        <template v-slot:cell(subject)="row">
                            <b-link @click="previewShare(row.item)">{{ row.item.subject }}</b-link>
                        </template>
                        <template v-slot:cell(created_at)="row">
                            <date-time :value="row.item.created_at"></date-time>
                        </template>
                        <template v-slot:cell(status)="row">
                            {{ statusText(row.item.status) }}
                        </template>
                        <template v-slot:cell(votes)="row">
                            <span v-if="row.item.status === 1">{{ row.item.like_count }}</span>
                            <span v-if="row.item.status > 1">{{ row.item.join_count }}</span>
                        </template>
                        <template v-slot:cell(opt)="row">
                             <b-button-group>
                                
                                <b-button size="sm" variant="success" v-if="canLike(row.item)" @click="iLikeIt(row.item.id, true)">感兴趣</b-button>
                                <b-button size="sm" variant="danger" v-if="canNotLike(row.item)" @click="iLikeIt(row.item.id, false)">不感兴趣</b-button>
                                <b-button size="sm" variant="success" v-if="row.item.status === 2 && !row.item.current_user_join && $store.getters.user.id != row.item.create_user_id" @click="iJoinIt(row.item.id, true)">参加</b-button>
                                <b-button size="sm" variant="danger" v-if="row.item.status === 2 && row.item.current_user_join && $store.getters.user.id != row.item.create_user_id" @click="iJoinIt(row.item.id, false)">不参加了</b-button>
                                
                                <b-button size="sm" variant="success" v-if="canPlan(row.item)" @click="createSharePlanDialog(row.item)">排期</b-button>
                                <b-button size="sm" variant="" v-if="canCancelPlan(row.item)" @click="cancelSharePlan(row.item)">取消排期</b-button>
                                <b-button size="sm" variant="warning" v-if="canFinishPlan(row.item)" @click="finishShareDialog(row.item)">完结</b-button>
                                <b-button size="sm" variant="danger" v-if="canDelete(row.item)">删除</b-button>

                                <b-button size="sm" variant="primary" v-if="row.item.status === 3" @click="previewShare(row.item)">回看</b-button>
                            </b-button-group>
                        </template>
                    </b-table>

                    <b-pagination v-model="current_page" :total-rows="page.total" :per-page="page.per_page"></b-pagination>
                </b-card-body>
            </b-card>

            <b-button v-b-modal.create-share-dialog>创建分享</b-button>
            <b-modal id="create-share-dialog" title="创建分享" hide-footer size="xl">
                <form @submit.stop.prevent="createShare">
                    <b-form-group id="subject-input-group" label="主题" label-for="subject-input">
                        <b-form-input id="subject-input" v-model="createForm.subject" type="text" placeholder="Enter subject" required></b-form-input>
                    </b-form-group>
                    <b-form-group id="subject-type-select-group" label="主题类型" label-for="subject-type-select">
                        <b-form-select v-model="createForm.subject_type" :options="subject_type_options"></b-form-select>
                    </b-form-group>
                    <b-form-group id="description-input-group" label="内容简介" label-for="description-input">
                        <b-form-textarea id="description-input" placeholder="Enter description" v-model="createForm.description" rows="6"/>
                    </b-form-group>

                    <b-form-group id="share_user-input-group" label="分享人" label-for="share_user-input" description="留空则分享人是自己">
                        <b-form-input id="share_user-input" v-model="createForm.share_user" type="text" placeholder="Enter share user"></b-form-input>
                    </b-form-group>
                    <b-button type="submit" variant="primary">创建</b-button>
                </form>
            </b-modal>
            <b-modal id="create-share-plan-dialog" :title="currentSharePlanTitle" hide-footer size="xl">
                <form @submit.stop.prevent="createSharePlan">
                    <b-form-group label="分享时间" >
                        <b-form-datepicker v-model="createPlanForm.share_date" class="mb-2"></b-form-datepicker>
                        <b-form-timepicker v-model="createPlanForm.share_time"></b-form-timepicker>
                    </b-form-group>
                    <b-form-group id="plan_duration-select-group" label="预计时长" label-for="plan_duration-select">
                        <b-form-select v-model="createPlanForm.plan_duration" :options="plan_duration_options"></b-form-select>
                    </b-form-group>
                    <b-form-group id="note-input-group" label="备注" label-for="note-input">
                        <b-form-textarea id="note-input" placeholder="Enter note" v-model="createPlanForm.note" rows="3"/>
                    </b-form-group>

                    <b-button type="submit" variant="primary">提交</b-button>
                </form>
            </b-modal>
            <b-modal id="finish-share-plan-dialog" :title="currentSharePlanTitle" hide-footer size="xl">
                <form @submit.stop.prevent="finishShare">
                    <b-form-group id="real_duration-select-group" label="实际时长" label-for="real_duration-select">
                        <b-form-select v-model="finishShareForm.real_duration" :options="plan_duration_options"></b-form-select>
                    </b-form-group>

                    <b-form-group id="attachment-upload-group" label="附件上传" label-for="attachment-upload-input">
                        <uploader :options="uploaderOptions" ref="uploader" :file-status-text="uploaderStatusText" @file-error="fileUploadFailed">
                            <uploader-unspport></uploader-unspport>
                            <uploader-drop>
                                <p>将文件拖拽到这里或者 </p>
                                <uploader-btn>选择要上传的文件</uploader-btn>
                            </uploader-drop>
                            <uploader-list></uploader-list>
                        </uploader>
                    </b-form-group>

                    <b-form-group id="note-input-group" label="备注" label-for="note-input">
                        <b-form-textarea id="note-input" placeholder="Enter note" v-model="finishShareForm.note" rows="3"/>
                    </b-form-group>

                    <b-button type="submit" variant="primary">提交</b-button>
                </form>
            </b-modal>
            <b-modal id="preview-share-dialog" :title="currentSharePlanTitle" hide-footer size="xl">
                <b-form>
                    <div v-if="shareDetail.share != null">
                        <b-form-group label="主题">
                            <b-form-input v-model="shareDetail.share.subject" type="text" readonly></b-form-input>
                        </b-form-group>
                        <b-form-group label="主题类型">
                            <b-form-select v-model="shareDetail.share.subject_type" :options="subject_type_options" disabled></b-form-select>
                        </b-form-group>
                        <b-form-group label="内容简介">
                            <b-form-textarea v-model="shareDetail.share.description" rows="4" readonly/>
                        </b-form-group>
                        <b-form-group label="分享人" >
                            <b-form-input v-model="shareDetail.share.share_user" type="text" readonly></b-form-input>
                        </b-form-group>
                    </div>

                    <div v-if="shareDetail.plan != null">
                        <b-form-group label="分享时间" readonly v-if="shareDetail.plan.share_at != ''">
                            <date-time :value="shareDetail.plan.share_at"></date-time>
                        </b-form-group>
                        <b-form-group label="预计时长" readonly v-if="shareDetail.plan.plan_duration > 0">
                            <b-form-input v-model="shareDetail.plan.plan_duration" type="text" readonly></b-form-input>
                        </b-form-group>
                        <b-form-group label="实际时长" readonly v-if="shareDetail.plan.real_duration > 0">
                            <b-form-input v-model="shareDetail.plan.real_duration" type="text" readonly></b-form-input>
                        </b-form-group>
                        <b-form-group label="备注">
                            <b-form-textarea v-model="shareDetail.plan.note" rows="4" readonly v-if="shareDetail.plan.note != ''"/>
                        </b-form-group>
                    </div>

                    <div v-if="shareDetail.attachments != null && shareDetail.attachments.length > 0 && shareDetail.plan != null">
                        <b-form-group label="附件">
                            <li v-for="(atta, index) in shareDetail.attachments" :key="index">
                                <a :href="'/storage/' + atta.atta_path" target="_blank">{{ atta.name }}</a>
                            </li>
                        </b-form-group>
                    </div>

                    <div v-if="current_share != null && current_share.create_user_id != this.$store.getters.user.id && current_share.status != 3">
                        <b-button size="sm" variant="success" v-if="canLike(current_share)" @click="iLikeItd(current_share.id, true)">感兴趣</b-button>
                        <b-button size="sm" variant="danger" v-if="canNotLike(current_share)" @click="iLikeItd(current_share.id, false)">不感兴趣</b-button>
                        <b-button size="sm" variant="success" v-if="current_share.status === 2 && !current_share.current_user_join" @click="iJoinItd(current_share.id, true)">参加</b-button>
                        <b-button size="sm" variant="danger" v-if="current_share.status === 2 && current_share.current_user_join" @click="iJoinItd(current_share.id, false)">不参加了</b-button>
                    </div>
                </b-form>
            </b-modal>
        </b-col>
    </b-row>
</template>

<script>
import axios from 'axios';
import moment from 'moment';

export default {
        name: 'Shares',
        components: {},
        data() {
            return {
                // 创建分享表单
                createForm: {
                    subject: '',
                    subject_type: '',
                    description: '',
                    share_user: this.$store.getters.user != null ? this.$store.getters.user.name : null,
                },
                // 创建分享计划表单
                createPlanForm: {
                    share_date: '',
                    share_time: '',
                    plan_duration: 0,
                },
                // 完成分享表单
                finishShareForm: {
                    real_duration: 0,
                    note: '',
                    attachments: [],
                },
                // 分享详情
                shareDetail: {share: {}, plan: {},},
                // 搜索表单
                search: {
                    status: this.QueryArgs(this.$route, 'status'),
                },
                // 分享状态选项列表
                status_options: [
                    {value: 0, text: '所有状态'},
                    {value: 1, text: '投票中'},
                    {value: 2, text: '已排期'},
                    {value: 3, text: '已完成'},
                ],
                // 分享主题类型选项列表
                subject_type_options: [
                    {value: '算法', text: '算法'},
                    {value: '机器学习', text: '机器学习'},
                    {value: '分布式系统', text: '分布式系统'},
                    {value: '数据库', text: '数据库'},
                ],
                // 预计分享时长选项列表
                plan_duration_options: [
                    {value: '10', text: '10 分钟'},
                    {value: '20', text: '20 分钟'},
                    {value: '30', text: '30 分钟'},
                    {value: '40', text: '40 分钟'},
                    {value: '50', text: '50 分钟'},
                    {value: '60', text: '60 分钟'},
                ],
                // 分享列表
                shares_fields: [
                    {key: 'id', label: 'ID'},
                    {key: 'subject', label: '主题'},
                    {key: 'subject_type', label: '分享类型'},
                    {key: 'share_user', label: '分享人'},
                    {key: 'votes', label: '票数'},
                    {key: 'status', label: '状态'},
                    {key: 'created_at', label: '创建时间'},
                    {key: 'opt', label: '操作'},
                ],
                shares: [],
                // 当前选中分享
                current_share: null,
                // 当前页码
                current_page: 1,
                // 分页信息
                page: {last_page: 1, page: 1, per_page: 20, total: 0},
                // 文件上传
                uploaderOptions: {
                    target: '/api/upload',
                    testChunks: false,
                },
                uploaderStatusText: {
                    success: '成功',
                    error: '错误',
                    uploading: '上传中',
                    paused: '暂停中',
                    waiting: '等待中',
                },
            };
        },
        computed: {
            currentSharePlanTitle() {
                if (this.current_share === null) {
                    return '';
                }

                return '分享计划：' + this.current_share.subject;
            },
        },
        watch: {
            '$route': 'reload',
            'current_page': 'reload',
        },
        methods: {
            // 文件上传失败
            fileUploadFailed(rootFile, file, message) {
                this.ToastError('文件上传失败：' + message);
            },
            // 分享搜索
            searchSubmit(evt) {
                evt.preventDefault();

                this.$router.push({path: '/', query: {
                    status: this.search.status,
                }}).catch(err => {this.ToastError(err);});
            },
            // 预览分享详情
            previewShare(share) {
                this.current_share = share;
                axios.get('/api/shares/' + share.id + '/').then(response => {
                    this.shareDetail = response.data;
                    this.$root.$emit('bv::show::modal', 'preview-share-dialog');
                }).catch(error => {this.ErrorBox(error)});
            },
            // 创建分享
            createShare() {
                axios.post('/api/shares/', this.createForm).then(response => {
                    this.$root.$emit('bv::hide::modal', 'create-share-dialog');
                    this.SuccessBox('创建成功');
                    this.reload();
                }).catch(error => {this.ErrorBox(error)})
            },
            // 创建分享计划
            createSharePlanDialog(share) {
                this.current_share = share;
                this.$root.$emit('bv::show::modal', 'create-share-plan-dialog');
            },
            createSharePlan() {
                let params = {};
                
                params.share_at = moment(this.createPlanForm.share_date + ' ' + this.createPlanForm.share_time, 'YYYY-MM-DD hh:mm:ss').format();
                params.plan_duration = parseInt(this.createPlanForm.plan_duration);
                params.note = this.createPlanForm.note;

                axios.post('/api/shares/' + this.current_share.id + '/plan', params).then(response => {
                    this.$root.$emit('bv::hide::modal', 'create-share-plan-dialog');
                    this.SuccessBox('操作成功');
                    this.reload();
                }).catch(error => {this.ErrorBox(error)});
            },
            cancelSharePlan(share) {
                axios.delete('/api/shares/' + share.id + '/plan').then(response => { 
                    this.SuccessBox('操作成功');
                    this.reload();
                }).catch(error => {this.ErrorBox(error)});
            },
            // 分享完结
            finishShareDialog(share) {
                this.current_share = share;
                this.$root.$emit('bv::show::modal', 'finish-share-plan-dialog');
            },
            finishShare() {
                let params = {};
                params.real_duration = parseInt(this.finishShareForm.real_duration);
                params.note = this.finishShareForm.note;
                params.attachments = [];
                for (let i in this.$refs.uploader.uploader.files) {
                   if (this.$refs.uploader.uploader.files[i].chunks[0] === undefined) { continue; }
                   params.attachments.push(JSON.parse(this.$refs.uploader.uploader.files[i].chunks[0].processedState.res).id);
                }

                axios.post('/api/shares/' + this.current_share.id + '/finish/', params).then(response => {
                    this.$root.$emit('bv::hide::modal', 'finish-share-plan-dialog');
                    this.SuccessBox('操作成功');
                    this.reload();
                }).catch(error => {this.ErrorBox(error)});
            },
            // 是否可以喜欢/参加
            canLike(share) {
                return share.status === 1 && !share.current_user_like && this.$store.getters.user.id != share.create_user_id;
            },
            canNotLike(share) {
                return share.status === 1 && share.current_user_like && this.$store.getters.user.id != share.create_user_id;
            },
            // 是否可以删除该分享
            canDelete(share) {
                return share.status === 1 && share.create_user_id == this.$store.getters.user.id;
            },
            // 是否可以排期
            canPlan(share) {
                return share.status === 1 && share.create_user_id == this.$store.getters.user.id;
            },
            canCancelPlan(share) {
                return share.status === 2 && share.create_user_id == this.$store.getters.user.id;
            },
            // 是否可以完结
            canFinishPlan(share) {
                return share.status === 2 && share.create_user_id == this.$store.getters.user.id;
            },
            // 喜欢or参加分享
            iLikeItd(shareId, like) {
                this.iLikeIt(shareId, like, () => {
                    this.$root.$emit('bv::hide::modal', 'preview-share-dialog');
                });
            },
            iJoinItd(shareId, join) {
                this.iJoinIt(shareId, join, () => {
                    this.$root.$emit('bv::hide::modal', 'preview-share-dialog');
                });
            },
            iLikeIt(shareId, like, success) {
                success = success || function() {};
                axios[like ? 'post':'delete']('/api/shares/' + shareId + '/like/').then(response => {
                    success();
                    this.ToastSuccess('操作成功');
                    this.reload();
                }).catch(error => {this.ErrorBox(error)})
            },
            iJoinIt(shareId, join, success) {
                success = success || function() {};
                axios[join ? 'post':'delete']('/api/shares/' + shareId + '/join/').then(response => {
                    success();
                    this.ToastSuccess('操作成功');
                    this.reload();
                }).catch(error => {this.ErrorBox(error)})
            },
            // 状态文本转换
            statusText(id) {
                switch (id) {
                case 1:
                    return '投票中';
                case 2: 
                    return '已排期';
                case 3: 
                    return '已完成';
                default: 
                    return '未知';
                }
            },
            // 页面刷新
            reload() {
                let params = this.$route.query;
                params.page = this.current_page;
                axios.get('/api/shares/', {
                    params: params,
                }).then(response => {
                    let shares = response.data.data;
                    let userLikeOrJoin = response.data.extra.user_like_or_join;
                    for (let i in shares) {
                        if (userLikeOrJoin[shares[i].id] === undefined) {
                            shares[i].current_user_like = false;
                            shares[i].current_user_join = false;
                        } else {
                            shares[i].current_user_like = userLikeOrJoin[shares[i].id].like;
                            shares[i].current_user_join = userLikeOrJoin[shares[i].id].join;
                        }
                    }

                    this.search = response.data.search;
                    this.shares = shares;
                    this.page = response.data.page;
                    this.current_page = this.page.page;

                    this.isBusy = false;
                }).catch(error => {
                    this.ToastError(error)
                });
            }
        },
        mounted() {
            this.reload();
        }
    }
</script>