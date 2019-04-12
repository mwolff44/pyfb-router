package db

// statements represents SQL queries to be prepared
var statements = map[string]string{
	"getBalance": `
		select pyfb_company.customer_balance, pyfb_customer.credit_limit, pyfb_customer.customer_enabled
		from pyfb_company
		join pyfb_customer on pyfb_customer.company_id = pyfb_company.id
		where pyfb_company.id=$1
	`,
	"getDirection": `
		select p.id as prefix_id, d.id, d.name
		from pyfb_direction_destination d
		join pyfb_direction_prefix p on d.id = p.destination_id
		AND $1 LIKE p.prefix || '%'
		order by length(p.prefix) desc limit 1;
	`,
	"getCustomerRate": `
		select * from (SELECT row_number() OVER () AS id,
		v.destnum_length_map,
		v.ratecard_id,
		v.rate_type,
		v.ratecard_name,
		v.rc_type,
		v.status,
		v.r_rate * 100000,
		v.r_block_min_duration,
		v.r_minimal_time,
		v.r_init_block * 100000,
		v.prefix,
		v.destnum_length,
		v.destination_id,
		v.country_id,
		v.type_id,
		v.region_id
			FROM ( select * from (select * from (SELECT r.id AS ratecard_id,
            1 AS rate_type,
            case 
              when (pr.destnum_length - length($2::text) = 0) then 0
              when (pr.destnum_length - length($2::text) <> 0 and pr.destnum_length = 0) then 1
              else 2 end as destnum_length_map,
            r.name AS ratecard_name,
            r.rc_type,
            pr.status,
            pr.r_rate,
            pr.r_block_min_duration,
            pr.r_minimal_time,
            pr.r_init_block,
            pr.prefix,
            pr.destnum_length,
            NULL::integer AS destination_id,
            NULL::integer AS country_id,
            NULL::integer AS type_id,
            NULL::integer AS region_id
           FROM pyfb_rating_c_ratecard r
             JOIN pyfb_rating_c_prefix_rate pr ON pr.c_ratecard_id = r.id AND pr.status::text <> 'disabled'::text and $2::text LIKE pr.prefix || '%'
          WHERE r.status::text = 'enabled'::text AND now() > r.date_start AND now() < r.date_end ) T WHERE destnum_length_map <> 2
          order by destnum_length_map, LENGTH(prefix) desc limit 1) E
        UNION ALL
         SELECT r.id AS ratecard_id,
            2 AS rate_type,
            0 as destnum_length_map,
            r.name AS ratecard_name,
            r.rc_type,
            dr.status,
            dr.r_rate,
            dr.r_block_min_duration,
            dr.r_minimal_time,
            dr.r_init_block,
            NULL::character varying AS prefix,
            0 AS destnum_length,
            dr.destination_id,
            NULL::integer AS country_id,
            NULL::integer AS type_id,
            NULL::integer AS region_id
           FROM pyfb_rating_c_ratecard r
             JOIN pyfb_rating_c_destination_rate dr ON dr.c_ratecard_id = r.id AND dr.status::text <> 'disabled'::text and dr.destination_id = $3
          WHERE r.status::text = 'enabled'::text AND now() > r.date_start AND now() < r.date_end
        UNION ALL
         SELECT r.id AS ratecard_id,
            3 AS rate_type,
            0 as destnum_length_map,
            r.name AS ratecard_name,
            r.rc_type,
            ctr.status,
            ctr.r_rate,
            ctr.r_block_min_duration,
            ctr.r_minimal_time,
            ctr.r_init_block,
            NULL::character varying AS prefix,
            0 AS destnum_length,
            NULL::integer AS destination_id,
            ctr.country_id,
            ctr.type_id,
            NULL::integer AS region_id
           FROM pyfb_rating_c_ratecard r
             join pyfb_direction_destination pdd on pdd.id = $3
						 join pyfb_direction_country pdc on pdd.country_iso2_id = pdc.country_iso2
						 join pyfb_direction_type pdt on pdd.type_id = pdt.id
						 JOIN pyfb_rating_c_countrytype_rate ctr ON ctr.c_ratecard_id = r.id AND ctr.status::text <> 'disabled'::text
						   and ctr.country_id = pdc.id and ctr.type_id = pdt.id
          WHERE r.status::text = 'enabled'::text AND now() > r.date_start AND now() < r.date_end
        UNION ALL
         SELECT r.id AS ratecard_id,
            4 AS rate_type,
            0 as destnum_length_map,
            r.name AS ratecard_name,
            r.rc_type,
            cr.status,
            cr.r_rate,
            cr.r_block_min_duration,
            cr.r_minimal_time,
            cr.r_init_block,
            NULL::character varying AS prefix,
            0 AS destnum_length,
            NULL::integer AS destination_id,
            cr.country_id,
            NULL::integer AS type_id,
            NULL::integer AS region_id
           FROM pyfb_rating_c_ratecard r
             join pyfb_direction_destination pdd on pdd.id = $3
             join pyfb_direction_country pdc on pdd.country_iso2_id = pdc.country_iso2
						 JOIN pyfb_rating_c_country_rate cr ON cr.c_ratecard_id = r.id AND cr.status::text <> 'disabled'::text
						   and cr.country_id = pdc.id
          WHERE r.status::text = 'enabled'::text AND now() > r.date_start AND now() < r.date_end
        UNION ALL
         SELECT r.id AS ratecard_id,
            5 AS rate_type,
            0 as destnum_length_map,
            r.name AS ratecard_name,
            r.rc_type,
            rtr.status,
            rtr.r_rate,
            rtr.r_block_min_duration,
            rtr.r_minimal_time,
            rtr.r_init_block,
            NULL::character varying AS prefix,
            0 AS destnum_length,
            NULL::integer AS destination_id,
            NULL::integer AS country_id,
            rtr.type_id,
            rtr.region_id
           FROM pyfb_rating_c_ratecard r
             join pyfb_direction_destination pdd on pdd.id = $3
						 join pyfb_direction_country pdc on pdd.country_iso2_id = pdc.country_iso2
						 join pyfb_direction_region pdr on pdc.region_id = pdr.id
             join pyfb_direction_type pdt on pdd.type_id = pdt.id
						 JOIN pyfb_rating_c_regiontype_rate rtr ON rtr.c_ratecard_id = r.id AND rtr.status::text <> 'disabled'::text
						   and rtr.region_id = pdr.id and rtr.type_id = pdt.id
          WHERE r.status::text = 'enabled'::text AND now() > r.date_start AND now() < r.date_end
        UNION ALL
         SELECT r.id AS ratecard_id,
            6 AS rate_type,
            0 as destnum_length_map,
            r.name AS ratecard_name,
            r.rc_type,
            rr.status,
            rr.r_rate,
            rr.r_block_min_duration,
            rr.r_minimal_time,
            rr.r_init_block,
            NULL::character varying AS prefix,
            0 AS destnum_length,
            NULL::integer AS destination_id,
            NULL::integer AS country_id,
            NULL::integer AS type_id,
            rr.region_id
           FROM pyfb_rating_c_ratecard r
             join pyfb_direction_destination pdd on pdd.id = $3
						 join pyfb_direction_country pdc on pdd.country_iso2_id = pdc.country_iso2
						 join pyfb_direction_region pdr on pdc.region_id = pdr.id
						 JOIN pyfb_rating_c_region_rate rr ON rr.c_ratecard_id = r.id AND rr.status::text <> 'disabled'::text
						   and rr.region_id = pdr.id
          WHERE r.status::text = 'enabled'::text AND now() > r.date_start AND now() < r.date_end
        UNION ALL
         SELECT r.id AS ratecard_id,
            7 AS rate_type,
            0 as destnum_length_map,
            r.name AS ratecard_name,
            r.rc_type,
            dr.status,
            dr.r_rate,
            dr.r_block_min_duration,
            dr.r_minimal_time,
            dr.r_init_block,
            NULL::character varying AS prefix,
            0 AS destnum_length,
            NULL::integer AS destination_id,
            NULL::integer AS country_id,
            NULL::integer AS type_id,
            NULL::integer AS region_id
           FROM pyfb_rating_c_ratecard r
             JOIN pyfb_rating_c_default_rate dr ON dr.c_ratecard_id = r.id AND dr.status::text <> 'disabled'::text
          WHERE r.status::text = 'enabled'::text AND now() > r.date_start AND now() < r.date_end) v)
          custrate where ratecard_id = $1 order by rate_type limit 1
		  ;
	`,
	/*si callerid_list_id not NULL, alors :
	-- - si destination_id est null => si callerid_filter = 3, on retient la ratecard (on a pas trouvé de prefixe interdit) sinon on ne la retient pas
	-- - si destination_id n'est pas null et si callerid_filter = 2 (auth) => OK si = 3 on ne retient pas la ratecard */
	"getCustomerRateCard": `
		with r as ( SELECT 
			cra.tech_prefix,
			cra.priority,
			cra.discount,
			cra.allow_negative_margin,
			cra.ratecard_id,
			c.name,
			c.callerid_list_id,
			prcl.callerid_filter, -- 2 authorized 3 unauthorized
			prcld.destination_id
			FROM pyfb_rating_cr_allocation cra
			  join pyfb_rating_c_ratecard c on c.id = cra.ratecard_id
				left join pyfb_rating_callernum_list prcl on prcl.id = c.callerid_list_id
				left join pyfb_rating_callernum_list_destination prcld on prcld.callernumlist_id = prcl.id and prcld.destination_id = $2
			WHERE c.status = 'enabled' AND cra.customer_id = $1 and c.rc_type = 'pstn' AND now() > c.date_start AND now() < c.date_end) 
		SELECT tech_prefix,
		  priority,
		  discount,
		  allow_negative_margin,
		  ratecard_id 
		FROM r WHERE r.callerid_list_id is null or 
		case when (r.callerid_list_id is not null and destination_id is null) then callerid_filter = '3' else case when (r.callerid_list_id is not null and destination_id is not null) then callerid_filter = '2' end end
		order by priority ASC;
	`,
	"getProviderRate": `
		select * from (SELECT row_number() OVER () AS id,
		v.destnum_length_map,
		v.ratecard_id,
		v.rate_type,
		v.ratecard_name,
		v.rc_type,
		v.status,
		v.r_rate * 100000,
		v.r_block_min_duration,
		v.r_minimal_time,
		v.r_init_block* 100000,
		v.prefix,
		v.destnum_length,
		v.destination_id,
		v.country_id,
		v.type_id,
		v.region_id
	FROM ( select * from (select * from (SELECT r.provider_id,
            r.id AS ratecard_id,
            case 
              when (pr.destnum_length - length($2::text) = 0) then 0
              when (pr.destnum_length - length($2::text) <> 0 and pr.destnum_length = 0) then 1
              else 2 end as destnum_length_map,
            1 AS rate_type,
            r.name AS ratecard_name,
            r.provider_prefix,
            r.estimated_quality,
            r.rc_type,
            pr.status,
            pr.r_rate,
            pr.r_block_min_duration,
            pr.r_minimal_time,
            pr.r_init_block,
            pr.prefix,
            pr.destnum_length,
            NULL::integer AS destination_id,
            NULL::integer AS country_id,
            NULL::integer AS type_id,
            NULL::integer AS region_id
           FROM pyfb_rating_p_ratecard r
             JOIN pyfb_rating_p_prefix_rate pr ON pr.p_ratecard_id = r.id AND pr.status::text <> 'disabled'::text and $2::text LIKE pr.prefix || '%'
          WHERE r.status::text = 'enabled'::text AND now() > r.date_start AND now() < r.date_end ) T WHERE destnum_length_map <> 2
          order by destnum_length_map, LENGTH(prefix) desc limit 1) E
        UNION ALL
         SELECT r.provider_id,
            r.id AS ratecard_id,
            0 as destnum_length_map,
            2 AS rate_type,
            r.name AS ratecard_name,
            r.provider_prefix,
            r.estimated_quality,
            r.rc_type,
            dr.status,
            dr.r_rate,
            dr.r_block_min_duration,
            dr.r_minimal_time,
            dr.r_init_block,
            NULL::character varying AS prefix,
            0 AS destnum_length,
            dr.destination_id,
            NULL::integer AS country_id,
            NULL::integer AS type_id,
            NULL::integer AS region_id
           FROM pyfb_rating_p_ratecard r
             JOIN pyfb_rating_p_destination_rate dr ON dr.p_ratecard_id = r.id AND dr.status::text <> 'disabled'::text and dr.destination_id = $3
          WHERE r.status::text = 'enabled'::text AND now() > r.date_start AND now() < r.date_end
        UNION ALL
         SELECT r.provider_id,
            r.id AS ratecard_id,
            0 as destnum_length_map,
            3 AS rate_type,
            r.name AS ratecard_name,
            r.provider_prefix,
            r.estimated_quality,
            r.rc_type,
            ctr.status,
            ctr.r_rate,
            ctr.r_block_min_duration,
            ctr.r_minimal_time,
            ctr.r_init_block,
            NULL::character varying AS prefix,
            0 AS destnum_length,
            NULL::integer AS destination_id,
            ctr.country_id,
            ctr.type_id,
            NULL::integer AS region_id
           FROM pyfb_rating_p_ratecard r
             join pyfb_direction_destination pdd on pdd.id = $3
						 join pyfb_direction_country pdc on pdd.country_iso2_id = pdc.country_iso2
						 join pyfb_direction_type pdt on pdd.type_id = pdt.id
						 JOIN pyfb_rating_p_countrytype_rate ctr ON ctr.p_ratecard_id = r.id AND ctr.status::text <> 'disabled'::text
						   and ctr.country_id = pdc.id and ctr.type_id = pdt.id
          WHERE r.status::text = 'enabled'::text AND now() > r.date_start AND now() < r.date_end
        UNION ALL
         SELECT r.provider_id,
            r.id AS ratecard_id,
            0 as destnum_length_map,
            4 AS rate_type,
            r.name AS ratecard_name,
            r.provider_prefix,
            r.estimated_quality,
            r.rc_type,
            cr.status,
            cr.r_rate,
            cr.r_block_min_duration,
            cr.r_minimal_time,
            cr.r_init_block,
            NULL::character varying AS prefix,
            0 AS destnum_length,
            NULL::integer AS destination_id,
            cr.country_id,
            NULL::integer AS type_id,
            NULL::integer AS region_id
           FROM pyfb_rating_p_ratecard r
             join pyfb_direction_destination pdd on pdd.id = $3
             join pyfb_direction_country pdc on pdd.country_iso2_id = pdc.country_iso2
						 JOIN pyfb_rating_p_country_rate cr ON cr.p_ratecard_id = r.id AND cr.status::text <> 'disabled'::text
						   and cr.country_id = pdc.id
          WHERE r.status::text = 'enabled'::text AND now() > r.date_start AND now() < r.date_end
        UNION ALL
         SELECT r.provider_id,
            r.id AS ratecard_id,
            0 as destnum_length_map,
            5 AS rate_type,
            r.name AS ratecard_name,
            r.provider_prefix,
            r.estimated_quality,
            r.rc_type,
            rtr.status,
            rtr.r_rate,
            rtr.r_block_min_duration,
            rtr.r_minimal_time,
            rtr.r_init_block,
            NULL::character varying AS prefix,
            0 AS destnum_length,
            NULL::integer AS destination_id,
            NULL::integer AS country_id,
            rtr.type_id,
            rtr.region_id
           FROM pyfb_rating_p_ratecard r
             join pyfb_direction_destination pdd on pdd.id = $3
						 join pyfb_direction_country pdc on pdd.country_iso2_id = pdc.country_iso2
						 join pyfb_direction_region pdr on pdc.region_id = pdr.id
             join pyfb_direction_type pdt on pdd.type_id = pdt.id
						 JOIN pyfb_rating_p_regiontype_rate rtr ON rtr.p_ratecard_id = r.id AND rtr.status::text <> 'disabled'::text
						   and rtr.region_id = pdr.id and rtr.type_id = pdt.id
          WHERE r.status::text = 'enabled'::text AND now() > r.date_start AND now() < r.date_end
        UNION ALL
         SELECT r.provider_id,
            r.id AS ratecard_id,
            0 as destnum_length_map,
            6 AS rate_type,
            r.name AS ratecard_name,
            r.provider_prefix,
            r.estimated_quality,
            r.rc_type,
            rr.status,
            rr.r_rate,
            rr.r_block_min_duration,
            rr.r_minimal_time,
            rr.r_init_block,
            NULL::character varying AS prefix,
            0 AS destnum_length,
            NULL::integer AS destination_id,
            NULL::integer AS country_id,
            NULL::integer AS type_id,
            rr.region_id
           FROM pyfb_rating_p_ratecard r
             join pyfb_direction_destination pdd on pdd.id = $3
						 join pyfb_direction_country pdc on pdd.country_iso2_id = pdc.country_iso2
						 join pyfb_direction_region pdr on pdc.region_id = pdr.id
						 JOIN pyfb_rating_p_region_rate rr ON rr.p_ratecard_id = r.id AND rr.status::text <> 'disabled'::text
						   and rr.region_id = pdr.id
          WHERE r.status::text = 'enabled'::text AND now() > r.date_start AND now() < r.date_end
        UNION ALL
         SELECT r.provider_id,
            r.id AS ratecard_id,
            0 as destnum_length_map,
            7 AS rate_type,
            r.name AS ratecard_name,
            r.provider_prefix,
            r.estimated_quality,
            r.rc_type,
            dr.status,
            dr.r_rate,
            dr.r_block_min_duration,
            dr.r_minimal_time,
            dr.r_init_block,
            NULL::character varying AS prefix,
            0 AS destnum_length,
            NULL::integer AS destination_id,
            NULL::integer AS country_id,
            NULL::integer AS type_id,
            NULL::integer AS region_id
           FROM pyfb_rating_p_ratecard r
             JOIN pyfb_rating_p_default_rate dr ON dr.p_ratecard_id = r.id AND dr.status::text <> 'disabled'::text
          WHERE r.status::text = 'enabled'::text AND now() > r.date_start AND now() < r.date_end) v)
  provrate where ratecard_id = $1 order by rate_type limit 1
          ;
	`,
	"getOutboundRoutes": `
	with r as (
		SELECT row_number() OVER () AS id,
			v.destnum_length_map,
			v.route_type,
			v.providerendpoint_id,
			v.provider_ratecard_id,
			v.route_rule,
			v.status,
			v.weight,
			v.priority,
			v.prefix,
			v.destnum_length,
			v.destination_id,
			v.country_id,
			v.type_id,
			v.region_id
		   FROM ( SELECT 
					case 
					  when (pr.destnum_length - length($2::text) = 0) then 0
					  when (pr.destnum_length - length($2::text) <> 0 and pr.destnum_length = 0) then 1
					  else 2 end as destnum_length_map,
					1 AS route_type,
					pg.providerendpoint_id,
					pr.provider_ratecard_id,
					pr.route_type AS route_rule,
					pr.status,
					pr.weight,
					pr.priority,
					pr.prefix,
					pr.destnum_length,
					NULL::integer AS destination_id,
					NULL::integer AS country_id,
					NULL::integer AS type_id,
					NULL::integer AS region_id
				   FROM pyfb_routing_c_routinggrp r
					 JOIN pyfb_routing_prefix_rule pr ON pr.c_route_id = r.routinggroup_id AND pr.status::text <> 'disabled'::text and $2::text LIKE pr.prefix || '%'
					 JOIN pyfb_routing_prefix_rule_provider_gateway_list pg ON pg.prefixrule_id = pr.id
				   where r.customer_id = $1
				UNION ALL
				 SELECT 
					0 as destnum_length_map,
					2 AS route_type,
					deg.providerendpoint_id,
					der.provider_ratecard_id,
					der.route_type AS route_rule,
					der.status,
					der.weight,
					der.priority,
					NULL::character varying AS prefix,
					0 AS destnum_length,
					der.destination_id,
					NULL::integer AS country_id,
					NULL::integer AS type_id,
					NULL::integer AS region_id
				   FROM pyfb_routing_c_routinggrp r
					 JOIN pyfb_routing_destination_rule der ON der.c_route_id = r.routinggroup_id AND der.status::text <> 'disabled'::text and der.destination_id = $3
					 JOIN pyfb_routing_destination_rule_provider_gateway_list deg ON deg.destinationrule_id = der.id
				   where r.customer_id = $1
				UNION ALL
				 SELECT 
					0 as destnum_length_map,
					3 AS route_type,
					ctg.providerendpoint_id,
					ctr.provider_ratecard_id,
					ctr.route_type AS route_rule,
					ctr.status,
					ctr.weight,
					ctr.priority,
					NULL::character varying AS prefix,
					0 AS destnum_length,
					NULL::integer AS destination_id,
					ctr.country_id,
					ctr.type_id,
					NULL::integer AS region_id
				   FROM pyfb_routing_c_routinggrp r
					 join pyfb_direction_destination pdd on pdd.id = $3
					 join pyfb_direction_country pdc on pdd.country_iso2_id = pdc.country_iso2
					 JOIN pyfb_routing_countrytype_rule ctr ON ctr.c_route_id = r.routinggroup_id AND ctr.status::text <> 'disabled'::text and ctr.country_id = pdc.id and ctr.type_id = pdd.type_id
					 JOIN pyfb_routing_countrytype_rule_provider_gateway_list ctg ON ctg.countrytyperule_id = ctr.id
				   where r.customer_id = $1
				UNION ALL
				 SELECT 
					0 as destnum_length_map,
					4 AS route_type,
					cg.providerendpoint_id,
					cr.provider_ratecard_id,
					cr.route_type AS route_rule,
					cr.status,
					cr.weight,
					cr.priority,
					NULL::character varying AS prefix,
					0 AS destnum_length,
					NULL::integer AS destination_id,
					cr.country_id,
					NULL::integer AS type_id,
					NULL::integer AS region_id
				   FROM pyfb_routing_c_routinggrp r
					 join pyfb_direction_destination pdd on pdd.id = $3
					 join pyfb_direction_country pdc on pdd.country_iso2_id = pdc.country_iso2
					 JOIN pyfb_routing_countryrule cr ON cr.c_route_id = r.routinggroup_id AND cr.status::text <> 'disabled'::text and cr.country_id = pdc.id
					 JOIN pyfb_routing_countryrule_provider_gateway_list cg ON cg.countryrule_id = cr.id
				   where r.customer_id = $1
				UNION ALL
				 SELECT 
					0 as destnum_length_map,
					5 AS route_type,
					rtg.providerendpoint_id,
					rtr.provider_ratecard_id,
					rtr.route_type AS route_rule,
					rtr.status,
					rtr.weight,
					rtr.priority,
					NULL::character varying AS prefix,
					0 AS destnum_length,
					NULL::integer AS destination_id,
					NULL::integer AS country_id,
					rtr.type_id,
					rtr.region_id
				   FROM pyfb_routing_c_routinggrp r
					 join pyfb_direction_destination pdd on pdd.id = $3
					 join pyfb_direction_country pdc on pdd.country_iso2_id = pdc.country_iso2
					 JOIN pyfb_routing_regiontype_rule rtr ON rtr.c_route_id = r.routinggroup_id AND rtr.status::text <> 'disabled'::text and rtr.region_id = pdc.region_id and rtr.type_id = pdd.type_id
					 JOIN pyfb_routing_regiontype_rule_provider_gateway_list rtg ON rtg.regiontyperule_id = rtr.id
				   where r.customer_id = $1
				UNION ALL
				 SELECT 
					0 as destnum_length_map,
					6 AS route_type,
					rg.providerendpoint_id,
					rr.provider_ratecard_id,
					rr.route_type AS route_rule,
					rr.status,
					rr.weight,
					rr.priority,
					NULL::character varying AS prefix,
					0 AS destnum_length,
					NULL::integer AS destination_id,
					NULL::integer AS country_id,
					NULL::integer AS type_id,
					rr.region_id
				   FROM pyfb_routing_c_routinggrp r
					 join pyfb_direction_destination pdd on pdd.id = $3
					 join pyfb_direction_country pdc on pdd.country_iso2_id = pdc.country_iso2
					 JOIN pyfb_routing_region_rule rr ON rr.c_route_id = r.routinggroup_id AND rr.status::text <> 'disabled'::text and rr.region_id = pdc.region_id
					 JOIN pyfb_routing_region_rule_provider_gateway_list rg ON rg.regionrule_id = rr.id
				   where r.customer_id = $1
				union all (
				 SELECT 
					0 as destnum_length_map,
					7 AS route_type,
					dg.providerendpoint_id,
					dr.provider_ratecard_id,
					dr.route_type AS route_rule,
					dr.status,
					dr.weight,
					dr.priority,
					NULL::character varying AS prefix,
					0 AS destnum_length,
					NULL::integer AS destination_id,
					NULL::integer AS country_id,
					NULL::integer AS type_id,
					NULL::integer AS region_id
				   FROM pyfb_routing_c_routinggrp r
					 JOIN pyfb_routing_default_rule dr ON dr.c_route_id = r.routinggroup_id AND dr.status::text <> 'disabled'::text
					 JOIN pyfb_routing_default_rule_provider_gateway_list dg ON dg.defaultrule_id = dr.id
					where r.customer_id = $1) 
				   ) v
				   )
		select 
		  r.*, 
			pep."name",
			pep.callee_norm_id,
			pep.callee_norm_in_id,
			pep.caller_id_in_from,
			pep.callerid_norm_id,
			pep.callerid_norm_in_id,
			pep.from_domain,
			pep.pai,
			pep.pid,
			pep.ppi,
			pep.prefix as gwprefix,
			pep.sip_transport,
			pep.sip_port,
			pep.sip_proxy,
			pep.username,
			pep.suffix,
			prpr.estimated_quality,
			prpr.provider_prefix as ratecard_prefix
		from r 
		join pyfb_endpoint_provider pep on r.providerendpoint_id = pep.id and pep.enabled = true 
		join pyfb_rating_p_ratecard prpr on pep.provider_id = prpr.provider_id and prpr.status = 'enabled'::text AND now() > prpr.date_start AND now() < prpr.date_end
		  where route_type = (select min(route_type) from r where destnum_length_map <> 2)
		  and case when route_type = 1 then destnum_length_map = (select min(destnum_length_map) from r where route_type = 1) else destnum_length_map = 0 end
		  and case when route_type = 1 then r.prefix = (select prefix from r where prefix is not null order by LENGTH(prefix) desc limit 1) else r.prefix is NULL end;
	`,
}
