<template>
    <b-row class="mb-5">
        <b-col>
            <b-card class="mb-2">
                <b-card-text>
                    <b-form inline @submit="searchSubmit">
                        <b-form-select v-model="search.status" class="mb-2 mr-sm-2 mb-sm-0" placeholder="状态" :options="status_options" v-if="action() !== 'recently'"></b-form-select>
                        <b-form-select v-model="search.type" class="mb-2 mr-sm-2 mb-sm-0" placeholder="类型" :options="subjectTypeFilterOptions"></b-form-select>
                        <b-button variant="light" type="submit">搜索</b-button>
                    </b-form>   
                </b-card-text>
            </b-card>

            <b-card class="mb-2" no-body v-if="shares.length > 0">
                <b-table :items="shares" :fields="shares_fields" :tbody-tr-class="sharesRowClass">
                    <template v-slot:cell(subject)="row">
                        <b-badge class="mr-2" :variant="randomStyle(row.item.subject_type)">{{ row.item.subject_type }}</b-badge>
                        <b-link v-b-tooltip.hover :title="'创建于 ' + dateFormat(row.item.created_at)" @click="previewShare(row.item)">{{ row.item.subject }}</b-link>
                        <b-link class="ml-2" target="_blank" :to="'/share?id=' + row.item.id"><font-awesome-icon icon="external-link-alt"></font-awesome-icon></b-link>
                    </template>
                    <template v-slot:cell(share_at)="row">
                        <span v-if="row.item.share_at != '0001-01-01T00:00:00Z'">
                            <date-time :value="row.item.share_at"></date-time>
                            <span v-if="row.item.share_room != ''"> / {{ row.item.share_room }}</span>
                            <span v-if="row.item.plan_duration > 0"> / {{ row.item.plan_duration }} 分钟</span>
                        </span>
                        <span v-else>-</span>
                    </template>
                    <template v-slot:cell(status)="row">
                        <span :class="statusText(row.item.status).class">{{ statusText(row.item.status).text }}</span> 
                        <span v-if="row.item.status != 3" class="ml-1" v-b-tooltip.hover :title="'当前 ' + voteCountSelector(row.item) + ' 人'">
                            <b-badge variant="" v-if="voteCountSelector(row.item) < 5">{{ voteCountSelector(row.item) }}</b-badge>
                            <b-badge variant="success" v-else-if="voteCountSelector(row.item) <= 10">{{ voteCountSelector(row.item) }}</b-badge>
                            <b-badge variant="warning" v-else>{{ voteCountSelector(row.item) }}</b-badge>
                        </span>
                    </template>

                    <template v-slot:cell(opt)="row">
                        <b-button-group>
                            
                            <b-button size="sm" variant="success" v-if="canLike(row.item)" @click="iLikeIt(row.item.id, true)">感兴趣</b-button>
                            <b-button size="sm" variant="danger" v-if="canNotLike(row.item)" @click="iLikeIt(row.item.id, false)">不感兴趣</b-button>
                            <b-button size="sm" variant="success" v-if="row.item.status === 2 && !row.item.current_user_join && $store.getters.user.id != row.item.create_user_id" @click="iJoinIt(row.item.id, true)">参加</b-button>
                            <b-button size="sm" variant="danger" v-if="row.item.status === 2 && row.item.current_user_join && $store.getters.user.id != row.item.create_user_id" @click="iJoinIt(row.item.id, false)">不参加了</b-button>
                            
                            <b-button size="sm" variant="info" v-if="canPlan(row.item)" @click="editShareDialog(row.item)">编辑</b-button>
                            <b-button size="sm" variant="success" v-if="canPlan(row.item)" @click="createSharePlanDialog(row.item)">排期</b-button>
                            <b-button size="sm" variant="info" v-if="canCancelPlan(row.item)" @click="editSharePlanDialog(row.item)">编辑排期</b-button>
                            <b-button size="sm" variant="warning" v-if="canFinishPlan(row.item)" @click="finishShareDialog(row.item)">完结</b-button>
                            <b-button size="sm" variant="warning" v-if="canEditFinishedPlan(row.item)" @click="finishShareDialog(row.item)">编辑</b-button>
                            <b-button size="sm" variant="danger" v-if="canDelete(row.item)" @click="deleteShare(row.item)">删除</b-button>

                            <b-button size="sm" variant="primary" v-if="row.item.status === 3" @click="previewShare(row.item)">回看</b-button>
                        </b-button-group>
                    </template>
                </b-table>
                <div class="mt-3" v-if="page.last_page > 1">
                    <b-pagination v-model="current_page" :total-rows="page.total" :per-page="page.per_page" align="center"></b-pagination>
                </div>
            </b-card>
            <b-card class="mb-2" v-if="shares.length == 0">当前没有相关分享</b-card>

            <b-button v-b-modal.create-share-dialog variant="primary">发起分享</b-button>
            <b-modal id="create-share-dialog" title="发起分享" hide-footer size="xl">
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

                    <b-form-group id="share_user-input-group" label="分享人" label-for="share_user-input">
                        <b-form-select v-model="createForm.share_user_id" :options="user_options"></b-form-select>
                        <b-form-input class="mt-2" id="share_user-input" v-model="createForm.share_user" type="text" placeholder="分享人姓名" v-if="createForm.share_user_id == 0"></b-form-input>
                    </b-form-group>
                    <b-button type="submit" variant="primary">确认发起</b-button>
                </form>
            </b-modal>
            <b-modal id="edit-share-dialog" title="编辑分享" hide-footer size="xl">
                <form @submit.stop.prevent="updateShare">
                    <b-form-group id="subject-input-group" label="主题" label-for="subject-input">
                        <b-form-input id="subject-input" v-model="createForm.subject" type="text" placeholder="Enter subject" required></b-form-input>
                    </b-form-group>
                    <b-form-group id="subject-type-select-group" label="主题类型" label-for="subject-type-select">
                        <b-form-select v-model="createForm.subject_type" :options="subject_type_options"></b-form-select>
                    </b-form-group>
                    <b-form-group id="description-input-group" label="内容简介" label-for="description-input">
                        <b-form-textarea id="description-input" placeholder="Enter description" v-model="createForm.description" rows="6"/>
                    </b-form-group>

                    <b-form-group id="share_user-input-group" label="分享人" label-for="share_user-input">
                        <b-form-select v-model="createForm.share_user_id" :options="user_options"></b-form-select>
                        <b-form-input class="mt-2" id="share_user-input" v-model="createForm.share_user" type="text" placeholder="分享人姓名" v-if="createForm.share_user_id == 0"></b-form-input>
                    </b-form-group>
                    <b-button type="submit" variant="primary">保存</b-button>
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
                    <b-form-group label="分享地点" label-for="share_room-input">
                        <b-form-input id="share_room-input" v-model="createPlanForm.share_room" type="text"></b-form-input>
                    </b-form-group>
                    <b-form-group id="note-input-group" label="备注" label-for="note-input">
                        <b-form-textarea id="note-input" placeholder="Enter note" v-model="createPlanForm.note" rows="3"/>
                    </b-form-group>

                    <b-button type="submit" variant="primary">保存</b-button>
                    <b-button class="ml-2" variant="danger" v-if="canCancelPlan(current_share)" @click="cancelSharePlan(current_share)">取消排期</b-button>
                </form>
            </b-modal>
            <b-modal id="finish-share-plan-dialog" :title="currentSharePlanTitle" hide-footer size="xl">
                <form @submit.stop.prevent="finishShare">
                    <b-form-group id="real_duration-select-group" label="实际时长" label-for="real_duration-select">
                        <b-form-select v-model="finishShareForm.real_duration" :options="plan_duration_options"></b-form-select>
                    </b-form-group>

                    <b-form-group label="已上传附件" v-if="finishShareForm.oldAttachments != null && finishShareForm.oldAttachments.length > 0">
                        <b-table :items="finishShareForm.oldAttachments" :fields="atta_fields">
                            <template v-slot:cell(opt)="row">
                                <b-button variant="danger" @click="deleteAttachments(row.index)">删除</b-button>
                            </template>
                            <template v-slot:cell(atta_name)="row">
                                <a :href="'/storage/' + row.item.atta_path" target="_blank">{{ row.item.name }}</a>
                            </template>
                        </b-table>
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
                        <b-form-group label="内容简介" v-if="shareDetail.share.description !== ''">
                            <b-form-textarea v-model="shareDetail.share.description" rows="4" readonly/>
                        </b-form-group>
                        <b-form-group label="分享人" v-if="shareDetail.share.share_user !== ''">
                            <b-form-input v-model="shareDetail.share.share_user" type="text" readonly></b-form-input>
                        </b-form-group>
                    </div>

                    <div v-if="(shareDetail.join_users !== null && shareDetail.join_users.length > 0) || (shareDetail.like_users !== null && shareDetail.like_users.length > 0)" class="mt-3">
                        <div v-if="shareDetail.like_users !== null && shareDetail.like_users.length > 0">
                            对分享感兴趣的人<br />
                            <b-badge class="mr-2" variant="info" v-for="(item, index) in shareDetail.like_users" :key="index">{{ item.name }}</b-badge>
                        </div>
                        <div v-if="shareDetail.join_users !== null && shareDetail.join_users.length > 0">
                            报名参加分享的人<br />
                            <b-badge class="mr-2" variant="primary" v-for="(item, index) in shareDetail.join_users" :key="index">{{ item.name }}</b-badge>
                        </div>
                    </div>

                    <div v-if="shareDetail.plan != null" class="mt-3">
                        <b-form-group label="分享时间" readonly v-if="shareDetail.plan.share_at !== ''">
                            <date-time :value="shareDetail.plan.share_at"></date-time>
                        </b-form-group>
                        <b-form-group label="分享地点" readonly v-if="shareDetail.plan.share_room !== ''">
                            <b-form-input v-model="shareDetail.plan.share_room" type="text" readonly></b-form-input>
                        </b-form-group>
                        <b-form-group label="预计时长" readonly v-if="shareDetail.plan.plan_duration > 0">
                            <b-form-input v-model="shareDetail.plan.plan_duration" type="text" readonly></b-form-input>
                        </b-form-group>
                        <b-form-group label="实际时长" readonly v-if="shareDetail.plan.real_duration > 0">
                            <b-form-input v-model="shareDetail.plan.real_duration" type="text" readonly></b-form-input>
                        </b-form-group>
                        <b-form-group label="备注" v-if="shareDetail.plan.note !== ''">
                            <b-form-textarea v-model="shareDetail.plan.note" rows="4" readonly />
                        </b-form-group>
                    </div>

                    <div v-if="shareDetail.attachments != null && shareDetail.attachments.length > 0 && shareDetail.plan != null" class="mt-3">
                        <b-form-group label="图片附件预览" v-if="imageFilter(shareDetail.attachments).length > 0">
                            <div v-for="(atta, index) in imageFilter(shareDetail.attachments)" :key="index" class="image-preview-box">
                                <img :src="'/storage/' + atta.atta_path" :title="atta.name">
                            </div>
                        </b-form-group>
                        <b-form-group label="附件">
                            <li v-for="(atta, index) in shareDetail.attachments" :key="index">
                                <a :href="'/storage/' + atta.atta_path" target="_blank">{{ atta.name }}</a>
                            </li>
                        </b-form-group>
                    </div>

                </b-form>
                <div v-if="current_share != null && current_share.create_user_id != this.$store.getters.user.id && current_share.status != 3" class="mt-3">
                    <b-button size="sm" variant="success" v-if="canLike(current_share)" @click="iLikeItd(current_share.id, true)">感兴趣</b-button>
                    <b-button size="sm" variant="danger" v-if="canNotLike(current_share)" @click="iLikeItd(current_share.id, false)">不感兴趣</b-button>
                    <b-button size="sm" variant="success" v-if="current_share.status === 2 && !current_share.current_user_join" @click="iJoinItd(current_share.id, true)">参加</b-button>
                    <b-button size="sm" variant="danger" v-if="current_share.status === 2 && current_share.current_user_join" @click="iJoinItd(current_share.id, false)">不参加了</b-button>
                </div>
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
                createForm: this.initCreateForm(),
                // 创建分享计划表单
                createPlanForm: this.initCreatePlanForm(),
                // 完成分享表单
                finishShareForm: this.initFinishShareForm(),
                // 分享详情
                shareDetail: this.initShareDetail(),
                // 搜索表单
                search: {
                    status: this.QueryArgs(this.$route, 'status'),
                    type: this.QueryArgs(this.$route, 'type'),
                },
                // 分享状态选项列表
                status_options: [
                    {value: 0, text: '所有状态'},
                    {value: 1, text: '意向收集'},
                    {value: 2, text: '报名'},
                    {value: 3, text: '已完成'},
                ],
                // 分享主题类型选项列表
                subject_type_options: [
                    {value: '算法&数据结构', text: '算法&数据结构'},
                    {value: '机器学习', text: '机器学习'},
                    {value: '分布式系统', text: '分布式系统'},
                    {value: '架构', text: '架构'},
                    {value: '数据库', text: '数据库'},
                    {value: '沟通的艺术', text: '沟通的艺术'},
                    {value: '职业规划', text: '职业规划'},
                    {value: '前端', text: '前端'},
                    {value: '编程语言', text: '编程语言'},
                    {value: '技术框架', text: '技术框架'},
                    {value: '产品', text: '产品'},
                    {value: 'UI设计&交互', text: 'UI设计&交互'},
                    {value: '大数据', text: '大数据'},
                    {value: '技术方案', text: '技术方案'},
                    {value: '其它', text: '其它'},
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
                    {key: 'subject', label: '主题'},
                    {key: 'share_user', label: '分享人'},
                    {key: 'status', label: '状态'},
                    {key: 'share_at', label: '分享计划'},
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
                // 随机样式map
                randomStyleMap: {},
                // 附件编辑表头
                atta_fields: [
                    {key: 'atta_name', label: '名称'},
                    {key: 'opt', label: '操作'},
                ],
                // 用户列表
                user_options: []
            };
        },
        computed: {
            currentSharePlanTitle() {
                if (this.current_share === null) {
                    return '';
                }

                return '分享计划：' + this.current_share.subject;
            },
            subjectTypeFilterOptions() {
                let items = [{value: '', text: '所有类型'}];
                items.push(...this.subject_type_options);
                return items;
            },
        },
        watch: {
            '$route': 'reload',
            'current_page': 'reload',
        },
        methods: {
            // 初始化
            initAllForm() {
                this.createForm = this.initCreateForm();
                this.createPlanForm = this.initCreatePlanForm();
                this.finishShareForm = this.initFinishShareForm();
                this.shareDetail = this.initShareDetail();
            },
            initCreateForm() {
                return {
                    subject: '',
                    subject_type: '',
                    description: '',
                    share_user: this.$store.getters.user != null ? this.$store.getters.user.name : null,
                    share_user_id: this.$store.getters.user != null ? this.$store.getters.user.id : null,
                }
            },
            initCreatePlanForm() {
                return {
                    share_date: '',
                    share_time: '',
                    share_room: '',
                    plan_duration: 0,
                    note: '',
                };
            },
            initFinishShareForm() {
                return {
                    real_duration: 0,
                    note: '',
                    oldAttachments: [],
                    attachments: [],
                };
            },
            initShareDetail() {
                return {share: {}, plan: {}, like_users: [], join_users: []};
            },
            // 当前页面动作类型
            action() {
                return this.QueryArgs(this.$route, 'act');
            },
            // 时间格式化
            dateFormat(d) {
                return moment(d).format('YYYY-MM-DD HH:mm:ss');
            },
            // 随机样式生成
            randomStyle(key) {
                if (this.randomStyleMap[key] === undefined) {
                    let items = ['primary', 'success', 'danger', 'warning', 'secondary', '', 'dark', 'light', 'info'];
                    this.randomStyleMap[key] = items[Math.floor(Math.random() * items.length)];
                }

                return this.randomStyleMap[key];
            },
            // 分享列表样式
            sharesRowClass(item) {
                if (item.status == 2) {
                    return 'table-success';
                }

                return '';
            },
            // 文件上传失败
            fileUploadFailed(rootFile, file, message) {
                this.ToastError('文件上传失败：' + message);
            },
            // 状态人数
            voteCountSelector(share) {
                switch (share.status) {
                    case 1:
                        return share.like_count;
                    case 2:
                        return share.join_count;
                    default:
                        return share.join_count;
                }
            },
            // 分享搜索
            searchSubmit(evt) {
                evt.preventDefault();

                let query = {};
                for (let i in this.$route.query) {
                    query[i] = this.$route.query[i];
                }
                for (let i in this.search) {
                    query[i] = this.search[i];
                }

                delete query.page;
                delete query.per_page;
                
                this.$router.push({path: '/', query: query}).catch(err => {this.ToastError(err);});
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
                axios.post('/api/shares/', this.createShareUpdateData()).then(() => {
                    this.$root.$emit('bv::hide::modal', 'create-share-dialog');
                    this.ToastSuccess('创建成功');
                    this.reload();
                }).catch(error => {this.ErrorBox(error)})
            },
            editShareDialog(share) {
                this.current_share = share;
                this.createForm.subject = share.subject;
                this.createForm.subject_type = share.subject_type;
                this.createForm.description = share.description;
                this.createForm.share_user = share.share_user;
                this.createForm.share_user_id = share.share_user_id;

                this.$root.$emit('bv::show::modal', 'edit-share-dialog');
            },
            createShareUpdateData() {
                let data = {}
                for (let k in this.createForm) {
                    data[k] = this.createForm[k];
                }

                if (this.createForm.share_user_id > 0) {
                    for (let i in this.user_options) {
                        if (this.user_options[i].value == this.createForm.share_user_id) {
                            data.share_user_id = this.user_options[i].value;
                            data.share_user = this.user_options[i].text;
                            break;
                        }
                    }  
                }

                return data;
            },
            updateShare() {
                axios.post('/api/shares/' + this.current_share.id + '/', this.createShareUpdateData()).then(() => {
                    this.$root.$emit('bv::hide::modal', 'edit-share-dialog');
                    this.ToastSuccess('修改成功');
                    this.reload();
                }).catch(error => {this.ErrorBox(error)})
            },
            // 创建分享计划
            createSharePlanDialog(share) {
                this.current_share = share;
                this.createPlanForm = this.initCreatePlanForm();
                this.$root.$emit('bv::show::modal', 'create-share-plan-dialog');
            },
            createSharePlan() {
                let params = {};
                
                params.share_at = moment(this.createPlanForm.share_date + ' ' + this.createPlanForm.share_time, 'YYYY-MM-DD hh:mm:ss').format();
                params.plan_duration = parseInt(this.createPlanForm.plan_duration);
                params.note = this.createPlanForm.note;
                params.share_room  = this.createPlanForm.share_room;

                axios.post('/api/shares/' + this.current_share.id + '/plan', params).then(() => {
                    this.$root.$emit('bv::hide::modal', 'create-share-plan-dialog');
                    this.ToastSuccess('操作成功');
                    this.reload();
                }).catch(error => {this.ErrorBox(error)});
            },
            editSharePlanDialog(share) {
                this.current_share = share;
                axios.get('/api/shares/' + share.id + '/').then(response => {
                    this.createPlanForm.share_date = moment(response.data.plan.share_at).format('YYYY-MM-DD');
                    this.createPlanForm.share_time = moment(response.data.plan.share_at).format('HH:mm');
                    this.createPlanForm.share_room = response.data.plan.share_room;
                    this.createPlanForm.plan_duration = response.data.plan.plan_duration;
                    this.createPlanForm.note = response.data.plan.note;

                    this.$root.$emit('bv::show::modal', 'create-share-plan-dialog');
                }).catch(error => {this.ErrorBox(error)});
            },
            cancelSharePlan(share) {
                this.$bvModal.msgBoxConfirm('确定取消该分享的排期？').then((value) => {
                    if (value !== true) {
                      return ;
                    }

                    axios.delete('/api/shares/' + share.id + '/plan').then(() => {
                        this.ToastSuccess('操作成功');
                        this.$root.$emit('bv::hide::modal', 'create-share-plan-dialog');
                        this.reload();
                    }).catch(error => {this.ErrorBox(error)});
                });
            },
            // 删除分享
            deleteShare(share) {
                this.$bvModal.msgBoxConfirm('确定删除该分享？').then((value) => {
                    if (value !== true) {
                        return ;
                    }

                    axios.delete('/api/shares/' + share.id + '/').then(() => {
                        this.ToastSuccess('操作成功');
                        this.reload();
                    });
                });
            },
            // 删除附件
            deleteAttachments(idx) {
                this.$bvModal.msgBoxConfirm('确定删除该附件？').then((value) => {
                    if (value !== true) {
                        return;
                    }

                    this.finishShareForm.oldAttachments.splice(idx, 1);
                });
            },
            // 分享完结
            finishShareDialog(share) {
                this.current_share = share;
                this.finishShareForm = this.initFinishShareForm();
                axios.get('/api/shares/' + share.id + '/').then(response => {
                    this.finishShareForm.real_duration = response.data.plan.real_duration;
                    this.finishShareForm.oldAttachments = response.data.attachments;
                    this.finishShareForm.attachments = response.data.share.attachments.split(",");
                    this.finishShareForm.note = response.data.plan.note;

                    this.$root.$emit('bv::show::modal', 'finish-share-plan-dialog');
                }).catch(error => {this.ErrorBox(error)});
                
            },
            finishShare() {
                let params = {};
                params.real_duration = parseInt(this.finishShareForm.real_duration);
                params.note = this.finishShareForm.note;
                params.attachments = [];

                for (let i in this.finishShareForm.oldAttachments) {
                    params.attachments.push(this.finishShareForm.oldAttachments[i].id);
                }

                for (let i in this.$refs.uploader.uploader.files) {
                   if (this.$refs.uploader.uploader.files[i].chunks[0] === undefined) { continue; }
                   params.attachments.push(JSON.parse(this.$refs.uploader.uploader.files[i].chunks[0].processedState.res).id);
                }

                axios.post('/api/shares/' + this.current_share.id + '/finish/', params).then(() => {
                    this.$root.$emit('bv::hide::modal', 'finish-share-plan-dialog');
                    this.ToastSuccess('操作成功');
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
                return share != null && share.status === 2 && share.create_user_id == this.$store.getters.user.id;
            },
            // 是否可以完结
            canFinishPlan(share) {
                return share.status === 2 && share.create_user_id == this.$store.getters.user.id;
            },
            // 是否可以编辑已经完结的分享
            canEditFinishedPlan(share) {
                return share.status === 3 && share.create_user_id == this.$store.getters.user.id;
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
                axios[like ? 'post':'delete']('/api/shares/' + shareId + '/like/').then(() => {
                    success();
                    this.ToastSuccess('操作成功');
                    this.reload();
                }).catch(error => {this.ErrorBox(error)})
            },
            iJoinIt(shareId, join, success) {
                success = success || function() {};
                axios[join ? 'post':'delete']('/api/shares/' + shareId + '/join/').then(() => {
                    success();
                    this.ToastSuccess('操作成功');
                    this.reload();
                }).catch(error => {this.ErrorBox(error)})
            },
            // 状态文本转换
            statusText(id) {
                switch (id) {
                case 1:
                    return {text: '意向收集', class: 'text-warning'};
                case 2: 
                    return {text: '报名', class: 'text-success'};
                case 3: 
                    return {text: '已完成', class: ''};
                default: 
                    return {text: '未知', class: ''};
                }
            },
            // 图片过滤
            imageFilter(attachments) {
                let results = [];
                for (let i in attachments) {
                    if (this.strIn(attachments[i].atta_type, ["jpg", "jpeg", "png", "gif", "bmp", "svg"])) {
                        results.push(attachments[i]);
                    }
                }

                return results;
            },
            strIn(str, arr) {
                for (let i in arr) {
                    if (arr[i].toLowerCase() === str.toLowerCase()) {
                        return true;
                    }
                }

                return false;
            },
            // 刷新所有用户选项
            reloadAllUsers() {
                axios.get('/api/users/').then(resp => {
                    let user_options = [{value: 0, text: '-- 外部联系人 --'}];
                    for (let i = 0; i < resp.data.length; i++) {
                        user_options.push({value: resp.data[i].id, text: resp.data[i].name})
                    }
                    
                    this.user_options = user_options;
                }).catch(error => {this.ToastError(error)});
            },
            // 页面刷新
            reload() {
                let sharesEndpoint = '/api/shares/';
                switch (this.action()) {
                case 'my':
                    sharesEndpoint = '/api/shares/my/';
                    break;
                case 'recently':
                    sharesEndpoint = '/api/shares/recently/';
                    break;
                }

                let params = this.$route.query;
                params.page = this.current_page;
                axios.get(sharesEndpoint, {
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
                }).catch(error => {
                    this.ToastError(error)
                });

                this.initAllForm();
                this.reloadAllUsers();
            }
        },
        mounted() {
            this.reload();
        }
    }
</script>

<style scoped>
.image-preview-box img {
    max-width: 600px;
    max-height: 600px;
    border: 10px solid #fff;
    margin-bottom: 10px;
}
</style>